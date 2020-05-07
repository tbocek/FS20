//go build auth.go storage.pb.go
package main

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/jessevdk/go-flags"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	optsService struct {
		Port int `short:"p" long:"port" description:"Port to listen on" required:"true"`
	}
	jwtKey = []byte("password123456")
)

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func auth(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		log.Printf("err %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if creds.Username != "user" || creds.Password != "password" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims := &Claims{
		Username: creds.Username, //this could be something app specific
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 5).Unix(),
			Id:        "15b54ed3-d1ed-4122-8e10-6c2f462ee398",
			Subject:   creds.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("token", tokenString)
	w.WriteHeader(http.StatusOK)
}

func main() {
	_, err := flags.NewParser(&optsService, flags.None).Parse()
	if err != nil {
		log.Fatal(err)
	}

	router := httprouter.New()
	router.POST("/auth", auth)

	log.Printf("Starting server on port %v...", optsService.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(optsService.Port), router))
}
