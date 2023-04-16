package routes

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net/http"
	"proxy/pkg/proxy/config"
)

// Root is the main handler for the proxy.
//
// It will forward the request to the service.
type Root struct {
	l  *zap.Logger
	c  *config.Proxy
	rc *redis.Client
	hc *http.Client
}

// NewRoot creates a new instance of Root.
func NewRoot(l *zap.Logger, c *config.Proxy, rc *redis.Client) *Root {
	return &Root{
		l:  l,
		c:  c,
		rc: rc,
		hc: &http.Client{
			Timeout: c.Request.Timeout,
		},
	}
}

func (r *Root) Proxy(w http.ResponseWriter, req *http.Request) {
	// check if the response is in the cache and return it.
	if req.Method == http.MethodGet {
		val, err := r.rc.Get(req.Context(), req.RequestURI).Result()

		switch err {
		case nil:
			r.l.Debug("cache hit", zap.String("url", req.RequestURI))
			sendResponse(w, val, http.StatusOK)
			return
		case redis.Nil:
			r.l.Debug("cache miss", zap.String("url", req.RequestURI))
		default:
			r.l.Error("cache error", zap.Error(err))
			sendResponse(w, err, http.StatusInternalServerError)
			return
		}
	}

	// configure the proxy request as a clone of the original request.
	pr, err := makeProxy(req, r.c.ServerURL)
	if err != nil {
		r.l.Error("proxy error", zap.Error(err))
		sendResponse(w, err, http.StatusInternalServerError)
		return
	}

	// forward the request to the service.
	resp, err := r.hc.Do(pr)
	if err != nil {
		r.l.Error("http error", zap.Error(err))
		sendResponse(w, err, http.StatusInternalServerError)
		return
	}

	// read the response body.
	body, err := readBody(resp)
	defer func() { _ = resp.Body.Close() }()
	if err != nil {
		r.l.Error("read error", zap.Error(err))
		sendResponse(w, err, http.StatusInternalServerError)
		return
	}

	// cache the response.
	if resp.StatusCode == http.StatusOK && req.Method == http.MethodGet {
		err = r.rc.Set(req.Context(), req.RequestURI, body, r.c.Cache.TTL).Err()

		if err != nil {
			r.l.Error("cache error", zap.Error(err))
		}
	}

	// send the response to the client.
	sendResponse(w, body, resp.StatusCode)
}
