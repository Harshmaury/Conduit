// @conduit-project: conduit
// @conduit-path: internal/session/registry_test.go
package session

import (
	"net"
	"testing"
	"time"
)

type fakeConn struct{ closed bool }

func (f *fakeConn) Close() error                       { f.closed = true; return nil }
func (f *fakeConn) Read(b []byte) (int, error)         { return 0, nil }
func (f *fakeConn) Write(b []byte) (int, error)        { return len(b), nil }
func (f *fakeConn) LocalAddr() net.Addr                { return nil }
func (f *fakeConn) RemoteAddr() net.Addr               { return nil }
func (f *fakeConn) SetDeadline(t time.Time) error      { return nil }
func (f *fakeConn) SetReadDeadline(t time.Time) error  { return nil }
func (f *fakeConn) SetWriteDeadline(t time.Time) error { return nil }

func makeSession(id, agentID, owner string) *Session {
	return &Session{
		ID: id, AgentID: agentID, OwnerSub: owner,
		Conn: &fakeConn{}, ConnectedAt: time.Now(), LastSeenAt: time.Now(),
	}
}

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := NewRegistry()
	s := makeSession("csn_1", "agent-a", "harsh@github")
	r.Register(s)
	if got := r.Get("csn_1"); got == nil || got.OwnerSub != "harsh@github" {
		t.Errorf("expected session, got %v", got)
	}
	if got := r.GetByAgent("agent-a"); got == nil {
		t.Error("expected session by agent")
	}
}

func TestRegistry_Remove(t *testing.T) {
	r := NewRegistry()
	conn := &fakeConn{}
	s := &Session{ID: "csn_2", AgentID: "agent-b", OwnerSub: "harsh@github",
		Conn: conn, ConnectedAt: time.Now(), LastSeenAt: time.Now()}
	r.Register(s)
	r.Remove("csn_2")
	if r.Get("csn_2") != nil {
		t.Error("expected nil after remove")
	}
	if !conn.closed {
		t.Error("expected connection closed")
	}
}

func TestRegistry_ReRegisterClosesOld(t *testing.T) {
	r := NewRegistry()
	old := &fakeConn{}
	r.Register(&Session{ID: "csn_3", AgentID: "agent-c", Conn: old,
		ConnectedAt: time.Now(), LastSeenAt: time.Now()})
	r.Register(&Session{ID: "csn_4", AgentID: "agent-c", Conn: &fakeConn{},
		ConnectedAt: time.Now(), LastSeenAt: time.Now()})
	if !old.closed {
		t.Error("expected old conn closed on re-register")
	}
	if got := r.GetByAgent("agent-c"); got == nil || got.ID != "csn_4" {
		t.Errorf("expected new session, got %v", got)
	}
}

func TestRegistry_Count(t *testing.T) {
	r := NewRegistry()
	if r.Count() != 0 {
		t.Error("expected 0")
	}
	r.Register(makeSession("csn_5", "a1", "u1"))
	r.Register(makeSession("csn_6", "a2", "u2"))
	if r.Count() != 2 {
		t.Errorf("expected 2, got %d", r.Count())
	}
}

func TestRegistry_Touch(t *testing.T) {
	r := NewRegistry()
	s := makeSession("csn_7", "a3", "u3")
	before := s.LastSeenAt
	time.Sleep(2 * time.Millisecond)
	r.Register(s)
	r.Touch("csn_7")
	got := r.Get("csn_7")
	if !got.LastSeenAt.After(before) {
		t.Error("expected LastSeenAt updated")
	}
}
