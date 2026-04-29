package server

import (
	"encoding/json"
	"net/http"

	"github.com/RyuichiroYoshida/imacan/backend/internal/auth"
	"github.com/RyuichiroYoshida/imacan/backend/internal/generated"
)

func NewRouter(handler *Handler, authService *auth.Service) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	strictHandler := generated.NewStrictHandlerWithOptions(
		handler,
		[]generated.StrictMiddlewareFunc{AuthMiddleware(authService)},
		generated.StrictHTTPServerOptions{
			RequestErrorHandlerFunc:  writeRequestError,
			ResponseErrorHandlerFunc: writeResponseError,
		},
	)

	return withCORS(generated.HandlerFromMux(strictHandler, mux))
}

func writeRequestError(w http.ResponseWriter, r *http.Request, err error) {
	writeJSON(w, http.StatusBadRequest, generated.ErrorResponseBody{
		Code:    "INVALID_REQUEST",
		Message: err.Error(),
	})
}

func writeResponseError(w http.ResponseWriter, r *http.Request, err error) {
	writeJSON(w, http.StatusInternalServerError, generated.ErrorResponseBody{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "internal server error",
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
