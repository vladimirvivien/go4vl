package main

import (
	"log"
	"net/http"
	"os"
)

var (
	port = ":5050"
)

func main() {
	if len(os.Args) > 2 {
		port = os.Args[1]
	}

	// serve examples dir
	log.Printf("serving files on port %s", port)
	http.Handle("/", http.FileServer(http.Dir("../")))
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
