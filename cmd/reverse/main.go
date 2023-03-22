package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"proxy/internal"
	"time"
)

type config struct {
	serverURL *url.URL
}

const (
	serverURL = "http://localhost:8081"
)

func initConfig() (*config, error) {
	s, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}

	return &config{serverURL: s}, nil
}

func configureRequest(r *http.Request, cfg *config) {
	r.Host = cfg.serverURL.Host
	r.URL.Host = cfg.serverURL.Host
	r.URL.Scheme = cfg.serverURL.Scheme
	r.RequestURI = ""
}

func initServerProxy(h http.Handler) *http.Server {
	return &http.Server{
		Addr:        ":8080",
		ReadTimeout: 5 * time.Second,
		Handler:     h,
	}
}

func main() {
	cfg, e := initConfig()
	if e != nil {
		panic(e)
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[INFO] request received at reverse proxy at %s\n", internal.GetCurrentTime())
		configureRequest(r, cfg)
		res, err := http.DefaultClient.Do(r)
		if err != nil {
			log.Printf("[ERROR] %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, err)
			return
		}

		defer func() { _ = res.Body.Close() }()
		log.Printf("[INFO] response received at reverse proxy at %s\n", internal.GetCurrentTime())
		w.WriteHeader(res.StatusCode)
		_, _ = io.Copy(w, res.Body)
	})

	s := initServerProxy(h)
	log.Println("[INFO] reverse proxy starting")
	log.Fatalf("[ERROR] %v", s.ListenAndServe())
}
