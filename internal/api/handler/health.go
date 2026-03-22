// @conduit-project: conduit
// @conduit-path: internal/api/handler/health.go
package handler

import (
	"encoding/json"
	"net/http"
)

// Health returns 200 OK — used by engx doctor (ADR-029).
func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": "conduit"})
}
