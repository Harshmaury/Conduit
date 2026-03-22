// @conduit-project: conduit
// @conduit-path: internal/session/eviction.go
// Eviction runs a background goroutine that removes ghost sessions.
// A ghost session is one where LastSeenAt > maxIdleTime — the remote
// engxa has disconnected silently without sending a close signal.
// Without this, crashed remote agents leave stale sessions in the registry
// indefinitely (Fix 2 — audit).
package session

import (
	"context"
	"log"
	"time"
)

const DefaultMaxIdleTime = 90 * time.Second
const evictionInterval  = 30 * time.Second

// StartEviction starts the background eviction loop.
// Runs until ctx is cancelled. Safe to call once from main.
func StartEviction(ctx context.Context, reg *Registry, maxIdle time.Duration, logger *log.Logger) {
	if maxIdle <= 0 {
		maxIdle = DefaultMaxIdleTime
	}
	go func() {
		ticker := time.NewTicker(evictionInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				evict(reg, maxIdle, logger)
			}
		}
	}()
}

// evict removes sessions idle longer than maxIdle.
func evict(reg *Registry, maxIdle time.Duration, logger *log.Logger) {
	cutoff := time.Now().Add(-maxIdle)
	sessions := reg.List()
	for _, s := range sessions {
		if s.LastSeenAt.Before(cutoff) {
			logger.Printf("evicting idle session id=%s agent=%s owner=%s idle=%.0fs",
				s.ID, s.AgentID, s.OwnerSub, time.Since(s.LastSeenAt).Seconds())
			reg.Remove(s.ID)
		}
	}
}
