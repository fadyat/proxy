package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Proxy is the configuration for the proxy.
type Proxy struct {

	// ProxyAddr is the address of the proxy.
	ProxyAddr string `envconfig:"PROXY_ADDR" required:"true"`

	// ServerURL is the URL of the server to proxy to.
	ServerURL string `envconfig:"PROXY_SERVER_URL" required:"true"`

	// RequestConfig is the configuration for the forwarded requests.
	Request Request

	// CacheConfig is the configuration for the cache.
	Cache Cache
}

// NewProxyConfig creates a new ProxyConfig.
//
// The configuration is loaded from environment variables, using envconfig.
// For local development, you can use a .env file and the godotenv package.
func NewProxyConfig() (*Proxy, error) {
	rc := &Request{}
	if err := envconfig.Process("request", rc); err != nil {
		return nil, err
	}

	cc := &Cache{}
	if err := envconfig.Process("cache", cc); err != nil {
		return nil, err
	}

	pc := &Proxy{
		Request: *rc, Cache: *cc,
	}
	if err := envconfig.Process("proxy", pc); err != nil {
		return nil, err
	}

	return pc, nil
}
