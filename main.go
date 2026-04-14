package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const rootfilepath = "."
	mux := http.NewServeMux() //Create ServeMux
	mux.Handle("/", http.FileServer(http.Dir(rootfilepath)))
	mux.Handle("/assets/logo.png", http.FileServer(http.Dir(rootfilepath)))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", rootfilepath, port)
	log.Fatal(srv.ListenAndServe())

}
