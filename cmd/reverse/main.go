package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"proxy/internal"
	"strconv"
	"time"
)

type config struct {
	serverURL *url.URL
	cache     *cacheCfg
}

type cacheCfg struct {
	addr     string
	password string
	db       int
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func getEnvOrDefaultInt(key string, defaultValue int) int {
	if value, ok := os.LookupEnv(key); ok {
		if v, err := strconv.Atoi(value); err == nil {
			return v
		}
	}

	return defaultValue
}

var (
	serverURL = getEnvOrDefault("SERVER_URL", "http://localhost:8081")
	redisAddr = getEnvOrDefault("REDIS_ADDR", "localhost:6379")
	redisPass = getEnvOrDefault("REDIS_PASS", "")
	redisDB   = getEnvOrDefaultInt("REDIS_DB", 0)
)

func initConfig() (*config, error) {
	s, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}

	return &config{
		serverURL: s,
		cache: &cacheCfg{
			addr:     redisAddr,
			password: redisPass,
			db:       redisDB,
		},
	}, nil
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

func root(w http.ResponseWriter, r *http.Request, cfg *config) {
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
}

func main() {
	cfg, err := initConfig()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		root(w, r, cfg)
	})

	cache := redis.NewClient(&redis.Options{
		Addr:     cfg.cache.addr,
		Password: cfg.cache.password,
		DB:       cfg.cache.db,
	})

	val, err := cache.Ping(context.TODO()).Result()
	if err != nil {
		panic(err)
	}
	log.Printf("[INFO] redis ping: %s", val)

	s := initServerProxy(mux)
	log.Println("[INFO] reverse proxy starting")
	log.Fatalf("[ERROR] %v", s.ListenAndServe())
}
