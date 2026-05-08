package workflow

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var (
	slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:_[a-z0-9]+)*$`)
)

type Config struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Version int    `json:"version"`
	Steps   []Step `json:"steps"`
}

type Step struct {
	Key          string     `json:"key"`
	Name         string     `json:"name"`
	ApproverRole string     `json:"approver_role"`
	Scope        string     `json:"scope"`
	Required     bool       `json:"required"`
	Condition    *Condition `json:"condition,omitempty"`
}

type Condition struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    any    `json:"value"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors"`
}

var builtInRoles = map[string]struct{}{
	"owner":       {},
	"admin":       {},
	"requester":   {},
	"manager":     {},
	"finance":     {},
	"procurement": {},
	"director":    {},
	"auditor":     {},
}

var supportedScopes = map[string]struct{}{
	"organization":         {},
	"requester_department": {},
}

var supportedOperators = map[string]struct{}{
	">":  {},
	">=": {},
	"<":  {},
	"<=": {},
	"==": {},
	"!=": {},
}

var supportedConditionFields = map[string]string{
	"estimated_amount": "number",
	"currency":         "string",
	"urgency":          "string",
	"department_id":    "string",
}

func ValidateConfig(cfg Config) ValidationResult {
	var errors []ValidationError

	if strings.TrimSpace(cfg.Name) == "" {
		errors = append(errors, ValidationError{Field: "name", Message: "name is required"})
	}

	if strings.TrimSpace(cfg.Type) != "purchase_request" {
		errors = append(errors, ValidationError{Field: "type", Message: "type must be purchase_request"})
	}

	if len(cfg.Steps) == 0 {
		errors = append(errors, ValidationError{Field: "steps", Message: "steps must contain at least one step"})
	}

	seenKeys := make(map[string]struct{}, len(cfg.Steps))
	for idx, step := range cfg.Steps {
		prefix := fmt.Sprintf("steps[%d]", idx)
		if strings.TrimSpace(step.Key) == "" {
			errors = append(errors, ValidationError{Field: prefix + ".key", Message: "step key is required"})
		} else if !slugPattern.MatchString(step.Key) {
			errors = append(errors, ValidationError{Field: prefix + ".key", Message: "step key must be slug-like"})
		} else if _, exists := seenKeys[step.Key]; exists {
			errors = append(errors, ValidationError{Field: prefix + ".key", Message: "step key must be unique"})
		} else {
			seenKeys[step.Key] = struct{}{}
		}

		if strings.TrimSpace(step.Name) == "" {
			errors = append(errors, ValidationError{Field: prefix + ".name", Message: "step name is required"})
		}

		if _, ok := builtInRoles[strings.TrimSpace(step.ApproverRole)]; !ok {
			errors = append(errors, ValidationError{Field: prefix + ".approver_role", Message: "approver_role must be a built-in role"})
		}

		if _, ok := supportedScopes[strings.TrimSpace(step.Scope)]; !ok {
			errors = append(errors, ValidationError{Field: prefix + ".scope", Message: "scope is not supported"})
		}

		if step.Condition != nil {
			errors = append(errors, validateCondition(prefix+".condition", *step.Condition)...)
		}
	}

	return ValidationResult{
		Valid:  len(errors) == 0,
		Errors: errors,
	}
}

func validateCondition(prefix string, condition Condition) []ValidationError {
	var errors []ValidationError

	expectedType, ok := supportedConditionFields[condition.Field]
	if !ok {
		errors = append(errors, ValidationError{Field: prefix + ".field", Message: "condition field is not supported"})
		return errors
	}

	if _, ok := supportedOperators[condition.Operator]; !ok {
		errors = append(errors, ValidationError{Field: prefix + ".operator", Message: "condition operator is not supported"})
	}

	switch expectedType {
	case "number":
		switch condition.Value.(type) {
		case float64, float32, int, int32, int64, uint, uint32, uint64, json.Number:
		default:
			errors = append(errors, ValidationError{Field: prefix + ".value", Message: "condition value must be numeric"})
		}
	case "string":
		if _, ok := condition.Value.(string); !ok {
			errors = append(errors, ValidationError{Field: prefix + ".value", Message: "condition value must be a string"})
		}
	}

	return errors
}

func EvaluateCondition(condition Condition, value any) (bool, error) {
	switch condition.Field {
	case "estimated_amount":
		left, ok := toFloat(value)
		if !ok {
			return false, fmt.Errorf("estimated_amount must be numeric")
		}
		right, ok := toFloat(condition.Value)
		if !ok {
			return false, fmt.Errorf("condition value must be numeric")
		}
		return compareNumbers(left, right, condition.Operator)
	case "currency", "urgency", "department_id":
		left, ok := value.(string)
		if !ok {
			return false, fmt.Errorf("%s must be a string", condition.Field)
		}
		right, ok := condition.Value.(string)
		if !ok {
			return false, fmt.Errorf("condition value must be a string")
		}
		return compareStrings(left, right, condition.Operator)
	default:
		return false, fmt.Errorf("condition field is not supported")
	}
}

func GenerateMermaid(cfg Config) string {
	var b strings.Builder
	b.WriteString("flowchart TD\n")
	b.WriteString("  A[Submit Purchase Request]")

	for i, step := range cfg.Steps {
		node := fmt.Sprintf("S%d", i+1)
		label := step.Name
		if step.Condition != nil {
			label = fmt.Sprintf("%s\\nif %s %s %v", label, step.Condition.Field, step.Condition.Operator, step.Condition.Value)
		}
		b.WriteString(fmt.Sprintf("\n  %s[%s]", node, escapeMermaid(label)))
		if i == 0 {
			b.WriteString(fmt.Sprintf("\n  A --> %s", node))
			continue
		}
		b.WriteString(fmt.Sprintf("\n  S%d --> %s", i, node))
	}

	if len(cfg.Steps) > 0 {
		b.WriteString(fmt.Sprintf("\n  S%d --> Z[Completed]", len(cfg.Steps)))
	}

	return b.String()
}

func escapeMermaid(text string) string {
	return strings.ReplaceAll(strings.ReplaceAll(text, "[", "("), "]", ")")
}

func toFloat(value any) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case json.Number:
		n, err := v.Float64()
		return n, err == nil
	default:
		return 0, false
	}
}

func compareNumbers(left, right float64, operator string) (bool, error) {
	switch operator {
	case ">":
		return left > right, nil
	case ">=":
		return left >= right, nil
	case "<":
		return left < right, nil
	case "<=":
		return left <= right, nil
	case "==":
		return left == right, nil
	case "!=":
		return left != right, nil
	default:
		return false, fmt.Errorf("condition operator is not supported")
	}
}

func compareStrings(left, right, operator string) (bool, error) {
	switch operator {
	case "==":
		return left == right, nil
	case "!=":
		return left != right, nil
	default:
		return false, fmt.Errorf("string conditions only support == and !=")
	}
}

func DefaultPurchaseRequestConfig() Config {
	return Config{
		Name:    "Purchase Request Approval",
		Type:    "purchase_request",
		Version: 1,
		Steps: []Step{
			{
				Key:          "manager_approval",
				Name:         "Manager Approval",
				ApproverRole: "manager",
				Scope:        "requester_department",
				Required:     true,
			},
			{
				Key:          "finance_approval",
				Name:         "Finance Approval",
				ApproverRole: "finance",
				Scope:        "organization",
				Required:     true,
				Condition: &Condition{
					Field:    "estimated_amount",
					Operator: ">",
					Value:    5000000,
				},
			},
			{
				Key:          "director_approval",
				Name:         "Director Approval",
				ApproverRole: "director",
				Scope:        "organization",
				Required:     true,
				Condition: &Condition{
					Field:    "estimated_amount",
					Operator: ">",
					Value:    25000000,
				},
			},
			{
				Key:          "procurement_processing",
				Name:         "Procurement Processing",
				ApproverRole: "procurement",
				Scope:        "organization",
				Required:     true,
			},
		},
	}
}
