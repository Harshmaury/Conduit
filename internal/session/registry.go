// @conduit-project: conduit
// @conduit-path: internal/session/registry.go
// Session registry — tracks active remote agent connections.
// ADR-042: every session carries a Gate-validated identity claim.
// Conduit is the remote agent bridge — each session is one engxa instance
// connecting from a remote machine.
package session

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// Session is one active remote agent connection.
type Session struct {
	ID           string    // UUID
	AgentID      string    // engxa agent identifier
	OwnerSub     string    // Gate subject (e.g. "harsh@github")
	Conn         net.Conn  // underlying connection from remote engxa
	ConnectedAt  time.Time
	LastSeenAt   time.Time
}

// Registry is a thread-safe in-memory store of active sessions.
type Registry struct {
	mu       sync.RWMutex
	byID     map[string]*Session
	byAgent  map[string]*Session // agentID → session (one session per agent)
}

// NewRegistry creates an empty Registry.
func NewRegistry() *Registry {
	return &Registry{
		byID:    make(map[string]*Session),
		byAgent: make(map[string]*Session),
	}
}

// Register adds a session. Replaces any existing session for the same agentID.
func (r *Registry) Register(s *Session) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if old, ok := r.byAgent[s.AgentID]; ok {
		old.Conn.Close()
		delete(r.byID, old.ID)
	}
	r.byID[s.ID] = s
	r.byAgent[s.AgentID] = s
}

// Remove removes a session by ID and closes its connection.
func (r *Registry) Remove(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	s, ok := r.byID[id]
	if !ok {
		return
	}
	s.Conn.Close()
	delete(r.byID, s.ID)
	delete(r.byAgent, s.AgentID)
}

// Get returns the session for the given ID, or nil.
func (r *Registry) Get(id string) *Session {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.byID[id]
}

// GetByAgent returns the session for the given agentID, or nil.
func (r *Registry) GetByAgent(agentID string) *Session {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.byAgent[agentID]
}

// List returns a snapshot of all active sessions.
func (r *Registry) List() []*Session {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*Session, 0, len(r.byID))
	for _, s := range r.byID {
		out = append(out, s)
	}
	return out
}

// Count returns the number of active sessions.
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.byID)
}

// Touch updates LastSeenAt for a session (heartbeat).
func (r *Registry) Touch(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if s, ok := r.byID[id]; ok {
		s.LastSeenAt = time.Now().UTC()
	}
}

// NewSessionID generates a unique session ID.
func NewSessionID() string {
	return fmt.Sprintf("csn_%x", time.Now().UnixNano())
}
