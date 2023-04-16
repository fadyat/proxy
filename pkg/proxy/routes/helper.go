package routes

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"proxy/pkg/proxy/responses"
)

// sendResponse sends a json response to the client.
func sendResponse(w http.ResponseWriter, val interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err, ok := val.(error)
	if ok {
		val = responses.Error{Message: err.Error()}
	}

	_ = json.NewEncoder(w).Encode(val)
}

// makeProxy creates a copy of the proxy request from the original request.
// Identical to the request passed to the service directly.
func makeProxy(r *http.Request, serviceURL string) (*http.Request, error) {
	clone := r.Clone(r.Context())

	u, err := url.Parse(serviceURL)
	if err != nil {
		return nil, err
	}

	clone.URL = u
	clone.Host = u.Host
	clone.RequestURI = ""

	return clone, nil
}

// readBody reads the body of the response as a string.
func readBody(r *http.Response) (string, error) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
