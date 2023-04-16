package config

import "time"

// Request contains the configuration for the http client.
type Request struct {

	// Timeout is the timeout for the request.
	Timeout time.Duration `envconfig:"REQUEST_TIMEOUT" required:"false" default:"5s"`
}
