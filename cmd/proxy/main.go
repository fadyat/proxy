package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net/http"
	"proxy/pkg"
	"proxy/pkg/proxy"
	"proxy/pkg/proxy/config"
	"proxy/pkg/proxy/routes"
	"time"
)

func main() {
	l := pkg.InitLogger()

	if err := godotenv.Load(".env"); err != nil {
		l.Warn("no .env file found")
	}

	cfg, err := config.NewProxyConfig()
	if err != nil {
		l.Panic("error while parsing config", zap.Error(err))
	}

	rc := redis.NewClient(&redis.Options{
		Addr:     cfg.Cache.RedisAddr,
		DB:       cfg.Cache.RedisDB,
		Password: cfg.Cache.RedisPass,
	})

	if err = rc.Ping(context.TODO()).Err(); err != nil {
		l.Panic("error while connecting to redis", zap.Error(err))
	}

	s := &http.Server{
		Addr:        cfg.ProxyAddr,
		ReadTimeout: 5 * time.Second,
		Handler: proxy.InitRoutes(
			routes.NewRoot(l, cfg, rc),
			l,
		),
	}

	l.Info("server starting", zap.String("addr", s.Addr))
	if err = s.ListenAndServe(); err != nil {
		l.Fatal("server failed", zap.Error(err))
	}
}
