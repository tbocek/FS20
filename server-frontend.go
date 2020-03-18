//go build server-frontend.go
package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
)

var (
	optsFrontend struct {
		Port int `short:"p" long:"port" description:"Port to listen on" required:"true"`
	}
)

func main() {
	_, err := flags.NewParser(&optsFrontend, flags.None).Parse()
	if err != nil {
		log.Fatal(err)
	}

	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte(fmt.Sprintf(`
<html>
 <head>
  <script type="text/javascript">
             fetch('http://127.0.0.1:8888/health')
               .then((response) => {
                 return response.json();
               })
               .then((data) => {
                 console.log(data);
               });
  </script>
 </head>
 <body>test</body>
</html>
		`)))
	})

	router.GET("/health", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.WriteHeader(http.StatusOK)
	})

	log.Printf("Starting server on port %v...", optsFrontend.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(optsFrontend.Port), router))
}
