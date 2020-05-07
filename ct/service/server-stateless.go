//go build auth.go storage.pb.go
package main

import (
	"context"
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
)

var (
	optsService struct {
		Port        int `short:"p" long:"port" description:"Port to listen on" required:"true"`
		StoragePort int `short:"s" long:"storage-port" description:"Storage port to connect to" required:"true"`
	}
	s      StorageClient
	jwtKey = []byte("password123456")
)

func getKey(w http.ResponseWriter, r *http.Request, claims *jwt.StandardClaims) {
	uuid := claims.Id
	v, err := s.GetKey(context.Background(), &Key{Key: uuid})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(v.Value) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.Copy(w, strings.NewReader(v.Value))
}

func postKey(w http.ResponseWriter, r *http.Request, claims *jwt.StandardClaims) {
	uuid := claims.Id
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("could not read request %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = s.PutKeyValue(context.Background(), &KeyValue{Key: uuid, Value: string(b)})
	if err != nil {
		log.Printf("could not store request %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func auth(next func(w http.ResponseWriter, r *http.Request, claims *jwt.StandardClaims)) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		authorizationHeader := req.Header.Get("Authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				claims := &jwt.StandardClaims{}
				token, err := jwt.ParseWithClaims(bearerToken[1], claims, func(token *jwt.Token) (i interface{}, err error) {
					return jwtKey, nil
				})
				if err != nil || !token.Valid {
					log.Printf("could not parse token %v", err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				next(w, req, claims)
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

	log.Println("Connecting to storage...")
	conn, err := grpc.Dial("172.17.0.1:"+strconv.Itoa(optsService.StoragePort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect to storage: %s", err)
	}
	defer conn.Close()
	s = NewStorageClient(conn)

	router := httprouter.New()
	router.GET("/portfolio", auth(getKey))
	router.POST("/portfolio", auth(postKey))

	log.Printf("Starting server on port %v...", optsService.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(optsService.Port), router))
}
