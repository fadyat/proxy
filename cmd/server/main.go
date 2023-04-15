package main

import (
	"fmt"
	"go.uber.org/zap"
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
	l := internal.InitLogger()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Info("received request", zap.String("url", r.RequestURI))
		_, _ = fmt.Fprintf(w, "Hello from server at %s", time.Now().Format(time.RFC3339))
		l.Info("response sent", zap.String("url", r.RequestURI))
	})

	s := initServer(h)
	l.Info("server starting", zap.String("addr", s.Addr))
	if err := s.ListenAndServe(); err != nil {
		l.Fatal("server failed", zap.Error(err))
	}
}
