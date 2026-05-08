package workflows

import (
	"testing"

	configworkflow "operra/api/internal/workflow"
)

func TestBuildConfigAppliesNameTypeAndVersion(t *testing.T) {
	svc := &Service{
		generate: func(cfg configworkflow.Config) string { return "diagram" },
		validate: configworkflow.ValidateConfig,
	}

	cfg, validation, err := svc.buildConfig("Workflow Name", "purchase_request", []byte(`{"steps":[{"key":"manager_approval","name":"Manager Approval","approver_role":"manager","scope":"requester_department","required":true}]}`), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !validation.Valid {
		t.Fatalf("expected valid config, got errors: %+v", validation.Errors)
	}
	if cfg.Name != "Workflow Name" {
		t.Fatalf("expected name overlay, got %q", cfg.Name)
	}
	if cfg.Type != "purchase_request" {
		t.Fatalf("expected type overlay, got %q", cfg.Type)
	}
	if cfg.Version != 1 {
		t.Fatalf("expected version 1, got %d", cfg.Version)
	}
}

func TestBuildConfigRejectsInvalidConfig(t *testing.T) {
	svc := &Service{
		generate: func(cfg configworkflow.Config) string { return "" },
		validate: configworkflow.ValidateConfig,
	}

	_, validation, err := svc.buildConfig("", "purchase_request", []byte(`{}`), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if validation.Valid {
		t.Fatal("expected invalid config")
	}
}
