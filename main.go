package main

import (
	"log"
	"net/http"
)

func main() {
	port := "8080"
	mux := http.NewServeMux()
	filepathRoot := http.Dir(".")
	mux.Handle("/", http.FileServer(filepathRoot))
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
