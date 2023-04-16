package config

import "time"

// Cache is the configuration for the redis cache.
type Cache struct {

	// RedisAddr is the address of the Redis server.
	RedisAddr string `envconfig:"CACHE_REDIS_ADDR" required:"true"`

	// RedisPass is the password of the Redis server.
	RedisPass string `envconfig:"CACHE_REDIS_PASS" required:"true"`

	// RedisDB is the database of the Redis server.
	RedisDB int `envconfig:"CACHE_REDIS_DB" required:"true"`

	// TTL is the time to live for the cache.
	TTL time.Duration `envconfig:"CACHE_TTL" required:"false" default:"1m"`
}
