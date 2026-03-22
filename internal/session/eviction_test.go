// @conduit-project: conduit
// @conduit-path: internal/session/eviction_test.go
package session

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestEvict_RemovesIdleSessions(t *testing.T) {
	r := NewRegistry()
	logger := log.New(os.Stderr, "", 0)

	// Add one fresh and one stale session
	fresh := makeSession("s1", "a1", "u1")
	stale := makeSession("s2", "a2", "u2")
	stale.LastSeenAt = time.Now().Add(-120 * time.Second) // 2 min ago

	r.Register(fresh)
	r.Register(stale)

	evict(r, 90*time.Second, logger)

	if r.Get("s1") == nil {
		t.Error("fresh session should not be evicted")
	}
	if r.Get("s2") != nil {
		t.Error("stale session should be evicted")
	}
}

func TestEvict_KeepsFreshSessions(t *testing.T) {
	r := NewRegistry()
	logger := log.New(os.Stderr, "", 0)

	s := makeSession("s3", "a3", "u3")
	r.Register(s)
	evict(r, 90*time.Second, logger)

	if r.Get("s3") == nil {
		t.Error("fresh session evicted incorrectly")
	}
}
