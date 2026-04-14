package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const rootfilepath = "."

	mux := http.NewServeMux() //Create ServeMux
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(rootfilepath))))
	// http.FileServer(http.Dir(rootfilepath))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", rootfilepath, port)
	log.Fatal(srv.ListenAndServe())

}
