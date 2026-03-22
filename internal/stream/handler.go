// @conduit-project: conduit
// @conduit-path: internal/stream/handler.go
// Package stream handles incoming remote engxa connections on port 9093.
// ADR-042: every connection must present a valid Gate identity token.
// A session is established only after identity is validated.
package stream

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/Harshmaury/Conduit/internal/auth"
	"github.com/Harshmaury/Conduit/internal/session"
)

// ConnectRequest is the JSON message a remote engxa sends on connection.
type ConnectRequest struct {
	AgentID       string `json:"agent_id"`
	IdentityToken string `json:"identity_token"` // Gate JWT (ADR-042)
}

// ConnectResponse is sent back after handshake.
type ConnectResponse struct {
	OK        bool   `json:"ok"`
	SessionID string `json:"session_id,omitempty"`
	Owner     string `json:"owner,omitempty"`
	Error     string `json:"error,omitempty"`
}

// Handler processes incoming remote agent connections.
type Handler struct {
	registry        *session.Registry
	validator       *auth.Validator
	requireIdentity bool
}

// NewHandler creates a stream Handler.
func NewHandler(reg *session.Registry, v *auth.Validator, requireIdentity bool) *Handler {
	return &Handler{registry: reg, validator: v, requireIdentity: requireIdentity}
}

// Handle runs the session lifecycle for one remote engxa connection.
// Blocks until the session closes. Safe to call in a goroutine.
func (h *Handler) Handle(conn net.Conn) {
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(10 * time.Second))
	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		sendError(conn, "handshake read failed")
		return
	}
	conn.SetDeadline(time.Time{})

	var req ConnectRequest
	if err := json.Unmarshal(scanner.Bytes(), &req); err != nil {
		sendError(conn, "handshake parse error")
		return
	}
	if req.AgentID == "" {
		sendError(conn, "agent_id is required")
		return
	}

	// Validate Gate identity (ADR-042).
	claim, err := h.validator.Validate(req.IdentityToken)
	if err != nil {
		sendError(conn, fmt.Sprintf("identity validation failed: %v", err))
		return
	}
	if h.requireIdentity && claim == nil {
		sendError(conn, "identity token required — authenticate with: engx login")
		return
	}

	ownerSub := ""
	if claim != nil {
		ownerSub = claim.Subject
		if !claim.HasScope("execute") {
			sendError(conn, "token missing required scope: execute")
			return
		}
	}

	sessionID := session.NewSessionID()
	s := &session.Session{
		ID:          sessionID,
		AgentID:     req.AgentID,
		OwnerSub:    ownerSub,
		Conn:        conn,
		ConnectedAt: time.Now().UTC(),
		LastSeenAt:  time.Now().UTC(),
	}
	h.registry.Register(s)
	defer h.registry.Remove(sessionID)

	resp := ConnectResponse{OK: true, SessionID: sessionID, Owner: ownerSub}
	b, _ := json.Marshal(resp)
	fmt.Fprintf(conn, "%s\n", b)

	// Keep session alive — heartbeat read loop.
	buf := make([]byte, 64)
	for {
		conn.SetDeadline(time.Now().Add(60 * time.Second))
		_, err := conn.Read(buf)
		if err != nil {
			break // client disconnected or timeout
		}
		h.registry.Touch(sessionID)
	}
}

func sendError(conn net.Conn, reason string) {
	resp := ConnectResponse{OK: false, Error: reason}
	b, _ := json.Marshal(resp)
	fmt.Fprintf(conn, "%s\n", b)
}
