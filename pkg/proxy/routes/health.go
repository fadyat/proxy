package routes

import (
	"net/http"
	"proxy/pkg/proxy/responses"
)

func Health(w http.ResponseWriter, _ *http.Request) {
	sendResponse(w, responses.Health{Status: "ok"}, http.StatusOK)
}
