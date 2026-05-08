package middleware

import "testing"

func TestHasAnyRole(t *testing.T) {
	user := &CurrentUser{Roles: []string{"requester", "finance"}}

	if !HasAnyRole(user, "admin", "finance") {
		t.Fatal("expected finance role to match")
	}

	if HasAnyRole(user, "owner", "admin") {
		t.Fatal("expected no match for roles user does not have")
	}
}
