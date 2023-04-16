package proxy

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"proxy/pkg/proxy/routes"
)

func logRequest(h http.Handler, l *zap.Logger) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Debug(
			"received request",
			zap.String("url", r.RequestURI),
			zap.String("method", r.Method),
		)

		h.ServeHTTP(w, r)
	})
}

func InitRoutes(rr *routes.Root, l *zap.Logger) http.Handler {
	mux := http.NewServeMux()

	// Health endpoint for the proxy.
	mux.HandleFunc("/q/health", func(w http.ResponseWriter, r *http.Request) {

		// If the request is coming to the proxy, we want to call the
		// health endpoint of the proxy.
		if r.Header.Get(HeaderXProxyProcess) != "true" {
			routes.Health(w, r)
			return
		}

		rr.Proxy(w, r)
	})

	// Metrics endpoint for the proxy.
	mux.HandleFunc("/q/metrics", promhttp.Handler().ServeHTTP)

	// Main proxy endpoint.
	mux.HandleFunc("/", rr.Proxy)

	// Log all requests.
	return logRequest(mux, l)
}
