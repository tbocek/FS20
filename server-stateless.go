//go build server-stateless.go storage.pb.go
package main

import (
	"context"
	"fmt"
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
	s StorageClient
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

func main() {
	_, err := flags.NewParser(&optsService, flags.None).Parse()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connecting to storage...")
	conn, err := grpc.Dial(":"+strconv.Itoa(optsService.StoragePort), grpc.WithInsecure())
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

	router.GET("/portfolio/:uuid", getKey)
	router.POST("/portfolio/:uuid", postKey)

	log.Printf("Starting server on port %v...", optsService.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(optsService.Port), router))
}
