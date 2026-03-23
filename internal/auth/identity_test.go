// @conduit-project: conduit
// @conduit-path: internal/auth/identity_test.go
package auth

import "testing"

func TestIdentityClaimDTO_HasScope(t *testing.T) {
	c := &IdentityClaimDTO{Scopes: []string{"execute", "observe"}}
	if !c.HasScope("execute") {
		t.Error("expected execute")
	}
	if c.HasScope("admin") {
		t.Error("should not have admin")
	}
	if c.HasScope("") {
		t.Error("empty scope should not match")
	}
}
