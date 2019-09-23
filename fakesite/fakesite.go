package fakesite

import (
	"fmt"
	"log"
	"net/http"
)

func Start(port string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Staing alive!")
	})
	log.Println("starting fake server on " + port)
	http.ListenAndServe(":"+port, nil)
}
