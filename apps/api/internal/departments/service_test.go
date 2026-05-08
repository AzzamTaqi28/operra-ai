package departments

import "testing"

func TestValidateDepartment(t *testing.T) {
	if err := validateDepartment(CreateRequest{Name: "Finance"}); err != nil {
		t.Fatalf("expected valid department, got %v", err)
	}

	if err := validateDepartment(CreateRequest{Name: ""}); err == nil {
		t.Fatal("expected missing name to fail")
	}
}
