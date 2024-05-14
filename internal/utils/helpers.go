package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/go-chi/chi/v5"
)

type DtoResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func SetupHostWithPort(host string, port int64) string { return fmt.Sprintf("%s:%d", host, port) }

func JsonResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(response)
	if err != nil {
		return
	}
}

func LogAllRoutes(r chi.Routes) {

	log := config.InitializeLogger().LogHandler

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		log.Debugf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		log.Errorf("Logging err: %s\n", err.Error())
	}
}
