package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	handler := &Handler{}
	uploadHandler := &UploadHandler{"localhost", "upload"}
	http.Handle("/", handler)
	http.Handle("/upload", uploadHandler)

	go func() {
		dirToSave := http.Dir(uploadHandler.UploadDir)
		fs := &http.Server{
			Addr:         ":8080",
			Handler:      http.FileServer(dirToSave),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		log.Fatal(fs.ListenAndServe())
	}()

	srv := &http.Server{
		Addr:         ":80",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
