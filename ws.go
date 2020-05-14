package main

import (
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	API_KEY = "57ddfdd2-c195-472c-8051-b321ff1612a2"
	URL     = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?start=1&limit=10&convert=USD"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func test(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "market-ws.html")
}

func ws(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	messageType, _, _ := conn.NextReader()
	for {
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
		w, _ := conn.NextWriter(messageType)
		io.Copy(w, res.Body)
		w.Close()
		time.Sleep(5 * time.Second)
	}
}

func main() {
	log.Print("Starting server...")
	router := httprouter.New()
	router.GET("/", test)
	router.GET("/ws", ws)
	log.Fatal(http.ListenAndServe(":8080", router))
}
