package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	sId  int32
	m    = make(map[string]string)
	opts struct {
		Port int `short:"p" long:"port" description:"Port to listen on" required:"true"`
	}
)

func get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uuid := ps.ByName("uuid")

	if len(m[uuid]) > 0 {
		w.Header().Set("Content-Type", "application/json")
		io.Copy(w, strings.NewReader(m[uuid]))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func post(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uuid := ps.ByName("uuid")

	if b, err := ioutil.ReadAll(r.Body); err == nil {
		m[uuid] = string(b)
		w.WriteHeader(http.StatusOK)
	} else {
		fmt.Printf("could not read request %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	sId = rand.Int31()

	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte(fmt.Sprintf("<html><body><h1>%v</h1></body></html>", sId)))
	})

	router.GET("/health", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.WriteHeader(http.StatusOK)
	})

	router.GET("/portfolio/:uuid", get)
	router.POST("/portfolio/:uuid", post)

	_, err := flags.NewParser(&opts, flags.None).Parse()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Starting server (%v) on port %v...", sId, opts.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(opts.Port), router))
}
