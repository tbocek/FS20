package main

import (
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
)

const (
	API_KEY = "57ddfdd2-c195-472c-8051-b321ff1612a2"
	URL     = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?start=1&limit=10&convert=USD"
)

func proxy(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		w.WriteHeader(400)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-CMC_PRO_API_KEY", API_KEY)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		w.WriteHeader(400)
	}
	w.Header().Set("Content-type", "application/json")
	io.Copy(w, res.Body)
}

func main() {
	log.Print("Starting server...")
	router := httprouter.New()
	router.GET("/", proxy)
	log.Fatal(http.ListenAndServe(":8080", router))
}
