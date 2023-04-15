package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"io"
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

func initServerProxy(h http.Handler) *http.Server {
	return &http.Server{
		Addr:        ":8080",
		ReadTimeout: 5 * time.Second,
		Handler:     h,
	}
}

func root(w http.ResponseWriter, r *http.Request, cfg *config, l *zap.Logger, cache *redis.Client) {
	l.Info("proxy request", zap.String("url", r.RequestURI))

	val, err := cache.Get(context.TODO(), r.RequestURI).Result()
	switch err {
	case redis.Nil:
		l.Debug("don't cached, doing request")
	case nil:
		w.WriteHeader(http.StatusOK)
		l.Info("cached, return cached value", zap.String("value", val))
		_, _ = fmt.Fprint(w, val)
		return
	default:
		l.Error("error while taking cached value", zap.Error(err))
	}

	r2 := r.Clone(context.Background())
	r2.URL = cfg.serverURL
	r2.Host = cfg.serverURL.Host
	r2.RequestURI = ""

	res, err := http.DefaultClient.Do(r2)
	if err != nil {
		zap.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, err)
		return
	}

	defer func() { _ = res.Body.Close() }()
	l.Info("proxy response", zap.String("status", res.Status))
	w.WriteHeader(res.StatusCode)
	_, _ = io.Copy(w, res.Body)

	serverBody := "value"
	if err != nil {
		l.Error("error while parsing body to json", zap.Error(err))
	}

	if err = cache.Set(context.TODO(), r.RequestURI, serverBody, time.Minute).Err(); err != nil {
		l.Error("error while caching value", zap.Error(err))
	} else {
		l.Debug("cached value", zap.String("key", r.RequestURI))
	}
}

func main() {
	cfg, err := initConfig()
	if err != nil {
		panic(err)
	}

	l := internal.InitLogger()
	cache := redis.NewClient(&redis.Options{
		Addr:     cfg.cache.addr,
		Password: cfg.cache.password,
		DB:       cfg.cache.db,
	})

	val, err := cache.Ping(context.TODO()).Result()
	if err != nil {
		panic(err)
	}
	l.Info("redis init", zap.String("value", val))

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		root(w, r, cfg, l, cache)
	})

	s := initServerProxy(mux)
	l.Info("server starting", zap.String("addr", s.Addr))
	if err = s.ListenAndServe(); err != nil {
		l.Fatal("server failed", zap.Error(err))
	}
}
