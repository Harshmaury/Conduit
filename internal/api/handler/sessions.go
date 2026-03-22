// @conduit-project: conduit
// @conduit-path: internal/api/handler/sessions.go
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Harshmaury/Conduit/internal/session"
)

// SessionDTO is the JSON representation of an active session.
type SessionDTO struct {
	ID          string `json:"id"`
	AgentID     string `json:"agent_id"`
	OwnerSub    string `json:"owner_sub"`
	ConnectedAt string `json:"connected_at"`
	LastSeenAt  string `json:"last_seen_at"`
}

// Sessions handles GET /conduit/sessions — returns all active remote sessions.
func Sessions(reg *session.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessions := reg.List()
		dtos := make([]SessionDTO, 0, len(sessions))
		for _, s := range sessions {
			dtos = append(dtos, SessionDTO{
				ID:          s.ID,
				AgentID:     s.AgentID,
				OwnerSub:    s.OwnerSub,
				ConnectedAt: s.ConnectedAt.Format("2006-01-02T15:04:05Z"),
				LastSeenAt:  s.LastSeenAt.Format("2006-01-02T15:04:05Z"),
			})
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"ok": true, "data": dtos})
	}
}
