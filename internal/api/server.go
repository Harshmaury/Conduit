// @conduit-project: conduit
// @conduit-path: internal/api/server.go
package api

import (
	"net/http"

	"github.com/Harshmaury/Conduit/internal/api/handler"
	"github.com/Harshmaury/Conduit/internal/config"
	"github.com/Harshmaury/Conduit/internal/session"
)

// NewServer builds the Conduit HTTP API server.
// Routes:
//   GET  /health              — liveness (always unauthenticated)
//   GET  /conduit/sessions    — list active remote sessions
func NewServer(reg *session.Registry, cfg *config.Config) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.Handle("/conduit/sessions", serviceAuthMiddleware(cfg.ServiceToken,
		handler.Sessions(reg)))
	return &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: mux,
	}
}

// serviceAuthMiddleware enforces X-Service-Token on protected endpoints (ADR-008).
func serviceAuthMiddleware(token string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token != "" && r.Header.Get(config.ServiceTokenHeader) != token {
			http.Error(w, `{"ok":false,"error":"UNAUTHORIZED"}`, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
