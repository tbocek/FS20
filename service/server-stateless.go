//go build server-stateless.go storage.pb.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jessevdk/go-flags"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	optsService struct {
		Port        int `short:"p" long:"port" description:"Port to listen on" required:"true"`
		StoragePort int `short:"s" long:"storage-port" description:"Storage port to connect to" required:"true"`
	}
	s      StorageClient
	jwtKey = []byte("password123456")
)

func getKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uuid := ps.ByName("uuid")
	w.Header().Set("Content-Type", "application/json")
	v, err := s.GetKey(context.Background(), &Key{Key: uuid})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(v.Value) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	io.Copy(w, strings.NewReader(v.Value))
}

func postKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uuid := ps.ByName("uuid")
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Printf("could not read request %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = s.PutKeyValue(context.Background(), &KeyValue{Key: uuid, Value: string(b)})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		fmt.Printf("err %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if creds.Username != "user" || creds.Password != "password" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
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

func auth(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		authorizationHeader := req.Header.Get("Authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (i interface{}, err error) {
					return jwtKey, nil
				})
				if err != nil || !token.Valid {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				next(w, req, ps)
				return
			}
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

func main() {
	_, err := flags.NewParser(&optsService, flags.None).Parse()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connecting to storage...")
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(optsService.StoragePort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect to storage: %s", err)
	}
	defer conn.Close()
	s = NewStorageClient(conn)

	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte(fmt.Sprintf("<html><body><h1>%v</h1></body></html>", optsService.Port)))
	})

	router.GET("/health", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.WriteHeader(http.StatusOK)
	})

	//router.POST("/login", login)
	//router.GET("/refresh", nil)

	router.GET("/portfolio/:uuid", auth(httprouter.Handle(getKey)))
	router.POST("/portfolio/:uuid", auth(httprouter.Handle(postKey)))

	log.Printf("Starting server on port %v...", optsService.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(optsService.Port), router))
}
