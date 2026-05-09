package ai

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"operra/api/internal/audit"
	configworkflow "operra/api/internal/workflow"
)

type Provider interface {
	Generate(ctx context.Context, prompt string) (*GeneratedWorkflow, error)
	Name() string
}

type GeneratedWorkflow struct {
	Config      configworkflow.Config `json:"config"`
	Explanation string                `json:"explanation"`
	Warnings    []string              `json:"warnings"`
	Model       string                `json:"model,omitempty"`
}

type Service struct {
	db       *sql.DB
	audit    *audit.Service
	provider Provider
}

func NewService(db *sql.DB, provider Provider, auditSvc ...*audit.Service) *Service {
	var svc *audit.Service
	if len(auditSvc) > 0 {
		svc = auditSvc[0]
	}
	return &Service{db: db, provider: provider, audit: svc}
}

type Result struct {
	WorkflowJSON   configworkflow.Config           `json:"workflow_json"`
	Explanation    string                          `json:"explanation"`
	MermaidDiagram string                          `json:"mermaid_diagram"`
	Validation     configworkflow.ValidationResult `json:"validation"`
	Warnings       []string                        `json:"warnings"`
	Provider       string                          `json:"provider"`
	Model          string                          `json:"model,omitempty"`
}

func (s *Service) Generate(ctx context.Context, organizationID, actorUserID, prompt string, validate func(configworkflow.Config) configworkflow.ValidationResult, mermaid func(configworkflow.Config) string) (*Result, error) {
	if strings.TrimSpace(prompt) == "" {
		return nil, errors.New("prompt is required")
	}
	if s.provider == nil {
		return nil, errors.New("ai provider is not configured")
	}

	generated, err := s.provider.Generate(ctx, prompt)
	if err != nil {
		return nil, err
	}

	validation := validate(generated.Config)
	result := &Result{
		WorkflowJSON:   generated.Config,
		Explanation:    generated.Explanation,
		MermaidDiagram: mermaid(generated.Config),
		Validation:     validation,
		Warnings:       generated.Warnings,
		Provider:       s.provider.Name(),
		Model:          generated.Model,
	}

	if err := s.logGeneration(ctx, organizationID, actorUserID, prompt, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) logGeneration(ctx context.Context, organizationID, actorUserID, prompt string, result *Result) error {
	output, err := json.Marshal(map[string]any{
		"workflow_json":   result.WorkflowJSON,
		"explanation":     result.Explanation,
		"mermaid_diagram": result.MermaidDiagram,
		"warnings":        result.Warnings,
		"provider":        result.Provider,
		"model":           result.Model,
	})
	if err != nil {
		return err
	}

	validationStatus := "valid"
	if !result.Validation.Valid {
		validationStatus = "invalid"
	}
	validationErrors, err := json.Marshal(result.Validation.Errors)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO ai_generation_logs (
			organization_id, actor_user_id, purpose, provider, model, input_prompt, output_json, validation_status, validation_errors
		) VALUES ($1, $2, 'workflow_generation', $3, $4, $5, $6, $7, $8)
	`, organizationID, actorUserID, result.Provider, result.Model, prompt, output, validationStatus, validationErrors)
	if err != nil {
		return err
	}

	if s.audit != nil {
		return s.audit.Record(ctx, organizationID, &actorUserID, "ai.workflow_generated", "ai_generation_log", nil, nil, map[string]any{
			"validation_status": validationStatus,
			"provider":          result.Provider,
		}, nil, nil)
	}

	return nil
}

type RuleBasedProvider struct{}

func NewRuleBasedProvider() *RuleBasedProvider {
	return &RuleBasedProvider{}
}

func (p *RuleBasedProvider) Name() string { return "rule_based" }

func (p *RuleBasedProvider) Generate(ctx context.Context, prompt string) (*GeneratedWorkflow, error) {
	cfg := configworkflow.DefaultPurchaseRequestConfig()
	cfg.Name = "Purchase Request Approval"

	promptLower := strings.ToLower(prompt)
	steps := make([]configworkflow.Step, 0, 4)
	steps = append(steps, managerStep())
	if strings.Contains(promptLower, "5 million") || strings.Contains(promptLower, "5000000") || strings.Contains(promptLower, "5,000,000") {
		steps = append(steps, financeStep(5000000))
	}
	if strings.Contains(promptLower, "25 million") || strings.Contains(promptLower, "25000000") || strings.Contains(promptLower, "25,000,000") {
		steps = append(steps, directorStep(25000000))
	}
	if !containsStep(steps, "procurement_processing") {
		steps = append(steps, procurementStep())
	}
	cfg.Steps = steps

	explanation := "This workflow uses role-based approval steps derived from the prompt."
	warnings := []string{}
	if !strings.Contains(promptLower, "procurement") {
		warnings = append(warnings, "procurement step added by default")
	}

	return &GeneratedWorkflow{
		Config:      cfg,
		Explanation: explanation,
		Warnings:    warnings,
		Model:       "rule-based",
	}, nil
}

type OpenAIProvider struct {
	baseURL string
	apiKey  string
	model   string
	client  *http.Client
}

func NewOpenAIProvider(baseURL, apiKey, model string) *OpenAIProvider {
	if strings.TrimSpace(model) == "" {
		model = "gpt-4.1-mini"
	}
	return &OpenAIProvider{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		model:   model,
		client:  &http.Client{Timeout: 45 * time.Second},
	}
}

func (p *OpenAIProvider) Name() string { return "openai" }

func (p *OpenAIProvider) Generate(ctx context.Context, prompt string) (*GeneratedWorkflow, error) {
	if strings.TrimSpace(p.apiKey) == "" {
		return nil, errors.New("openai api key is required")
	}

	payload := map[string]any{
		"model": p.model,
		"messages": []map[string]string{
			{"role": "system", "content": "Return only valid JSON with keys config, explanation, warnings. config must be a purchase_request workflow JSON."},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.2,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ai provider error: %s", resp.Status)
	}

	var envelope struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, err
	}
	if len(envelope.Choices) == 0 {
		return nil, errors.New("empty ai response")
	}

	var parsed struct {
		Config      configworkflow.Config `json:"config"`
		Explanation string                `json:"explanation"`
		Warnings    []string              `json:"warnings"`
	}
	if err := json.Unmarshal([]byte(envelope.Choices[0].Message.Content), &parsed); err != nil {
		return nil, err
	}

	return &GeneratedWorkflow{
		Config:      parsed.Config,
		Explanation: parsed.Explanation,
		Warnings:    parsed.Warnings,
		Model:       p.model,
	}, nil
}

func managerStep() configworkflow.Step {
	return configworkflow.Step{Key: "manager_approval", Name: "Manager Approval", ApproverRole: "manager", Scope: "requester_department", Required: true}
}

func financeStep(threshold float64) configworkflow.Step {
	return configworkflow.Step{
		Key:          "finance_approval",
		Name:         "Finance Approval",
		ApproverRole: "finance",
		Scope:        "organization",
		Required:     true,
		Condition:    &configworkflow.Condition{Field: "estimated_amount", Operator: ">", Value: threshold},
	}
}

func directorStep(threshold float64) configworkflow.Step {
	return configworkflow.Step{
		Key:          "director_approval",
		Name:         "Director Approval",
		ApproverRole: "director",
		Scope:        "organization",
		Required:     true,
		Condition:    &configworkflow.Condition{Field: "estimated_amount", Operator: ">", Value: threshold},
	}
}

func procurementStep() configworkflow.Step {
	return configworkflow.Step{Key: "procurement_processing", Name: "Procurement Processing", ApproverRole: "procurement", Scope: "organization", Required: true}
}

func containsStep(steps []configworkflow.Step, key string) bool {
	for _, step := range steps {
		if step.Key == key {
			return true
		}
	}
	return false
}
