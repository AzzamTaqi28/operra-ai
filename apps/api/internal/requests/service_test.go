package requests

import (
	"testing"

	configworkflow "operra/api/internal/workflow"
)

func TestValidateCreateRequest(t *testing.T) {
	req := CreateRequest{
		DepartmentID:    "dept-1",
		Title:           "Buy laptops",
		ItemName:        "Laptop",
		Description:     "Need devices",
		Quantity:        2,
		EstimatedAmount: 1000000,
	}
	if err := validateCreateRequest(req); err != nil {
		t.Fatalf("expected request to be valid, got %v", err)
	}

	if err := validateCreateRequest(CreateRequest{}); err == nil {
		t.Fatal("expected empty request to fail")
	}
}

func TestParseDate(t *testing.T) {
	parsed, err := parseDate("2026-07-01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if parsed == nil {
		t.Fatal("expected parsed date")
	}

	if _, err := parseDate("07/01/2026"); err == nil {
		t.Fatal("expected invalid date format to fail")
	}
}

func TestEvaluateStepCondition(t *testing.T) {
	request := &PurchaseRequest{
		OrganizationID:  "org-1",
		ID:              "req-1",
		DepartmentID:    "dept-1",
		EstimatedAmount: 10000000,
		Currency:        "IDR",
		Urgency:         "urgent",
	}

	ok, err := evaluateStepCondition(configworkflow.Step{
		Condition: &configworkflow.Condition{
			Field:    "estimated_amount",
			Operator: ">",
			Value:    5000000,
		},
	}, request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected amount condition to pass")
	}

	ok, err = evaluateStepCondition(configworkflow.Step{
		Condition: &configworkflow.Condition{
			Field:    "urgency",
			Operator: "==",
			Value:    "urgent",
		},
	}, request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected urgency condition to pass")
	}
}

func TestApplicableStepsThresholds(t *testing.T) {
	cfg := configworkflow.DefaultPurchaseRequestConfig()

	steps, err := applicableSteps(cfg, &PurchaseRequest{EstimatedAmount: 1000000, Currency: "IDR", Urgency: "normal"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(steps) != 2 {
		t.Fatalf("expected 2 steps for 1,000,000 amount, got %d", len(steps))
	}

	steps, err = applicableSteps(cfg, &PurchaseRequest{EstimatedAmount: 10000000, Currency: "IDR", Urgency: "normal"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(steps) != 3 {
		t.Fatalf("expected 3 steps for 10,000,000 amount, got %d", len(steps))
	}

	steps, err = applicableSteps(cfg, &PurchaseRequest{EstimatedAmount: 30000000, Currency: "IDR", Urgency: "normal"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(steps) != 4 {
		t.Fatalf("expected 4 steps for 30,000,000 amount, got %d", len(steps))
	}
}
