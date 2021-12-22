package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const fsPort = 8089

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/upload", uploadHandler)
	mux.HandleFunc("/list", listHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", fsPort),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Printf("api start listening on %d port\n", fsPort)
	log.Fatal(srv.ListenAndServe())
}
