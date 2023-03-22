package main

import (
	"fmt"
	"log"
	"net/http"
	"proxy/internal"
	"time"
)

func initServer(h http.Handler) *http.Server {
	return &http.Server{
		Addr:        ":8081",
		ReadTimeout: 5 * time.Second,
		Handler:     h,
	}
}

func main() {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[INFO] request received at server at %s", internal.GetCurrentTime())
		_, _ = fmt.Fprintf(w, "Hello from server at %s", internal.GetCurrentTime())
		log.Printf("[INFO] response sent at server at %s", internal.GetCurrentTime())
	})

	s := initServer(h)
	log.Println("[INFO] server starting")
	log.Fatalf("[ERROR] %v", s.ListenAndServe())
}
