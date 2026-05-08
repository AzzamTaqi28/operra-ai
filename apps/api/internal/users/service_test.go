package users

import "testing"

func TestNormalizeRoleKeys(t *testing.T) {
	keys := normalizeRoleKeys([]string{"manager", "manager", "finance", " "})
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %v", keys)
	}
	if keys[0] != "manager" || keys[1] != "finance" {
		t.Fatalf("unexpected normalized keys: %v", keys)
	}
}

func TestValidateCreate(t *testing.T) {
	if err := validateCreate(CreateRequest{Name: "A", Email: "a@example.com", Password: "secret"}); err != nil {
		t.Fatalf("expected request to be valid, got %v", err)
	}

	if err := validateCreate(CreateRequest{}); err == nil {
		t.Fatal("expected empty request to fail")
	}
}
