package workflow

import "testing"

func TestValidateConfig(t *testing.T) {
	result := ValidateConfig(DefaultPurchaseRequestConfig())
	if !result.Valid {
		t.Fatalf("expected config to be valid, got errors: %+v", result.Errors)
	}
}

func TestValidateConfigRejectsBadConfig(t *testing.T) {
	result := ValidateConfig(Config{})
	if result.Valid {
		t.Fatal("expected config to be invalid")
	}
}

func TestEvaluateCondition(t *testing.T) {
	ok, err := EvaluateCondition(Condition{Field: "estimated_amount", Operator: ">", Value: 5000000}, 10000000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected condition to be true")
	}

	ok, err = EvaluateCondition(Condition{Field: "urgency", Operator: "==", Value: "urgent"}, "urgent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected string condition to be true")
	}
}

func TestGenerateMermaid(t *testing.T) {
	diagram := GenerateMermaid(DefaultPurchaseRequestConfig())
	if diagram == "" {
		t.Fatal("expected mermaid output")
	}
}
