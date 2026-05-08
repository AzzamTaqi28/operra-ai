package requests

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	configworkflow "operra/api/internal/workflow"
)

type WorkflowLookup struct {
	WorkflowID string
	VersionID  string
	Config     configworkflow.Config
}

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

type PurchaseRequest struct {
	ID                    string                 `json:"id"`
	OrganizationID        string                 `json:"organization_id"`
	WorkflowID            *string                `json:"workflow_id,omitempty"`
	WorkflowVersionID     *string                `json:"workflow_version_id,omitempty"`
	RequesterID           string                 `json:"requester_id"`
	DepartmentID          string                 `json:"department_id"`
	Title                 string                 `json:"title"`
	ItemName              string                 `json:"item_name"`
	Description           string                 `json:"description"`
	Quantity              float64                `json:"quantity"`
	EstimatedAmount       float64                `json:"estimated_amount"`
	Currency              string                 `json:"currency"`
	Urgency               string                 `json:"urgency"`
	ExpectedDate          *string                `json:"expected_date,omitempty"`
	VendorName            *string                `json:"vendor_name,omitempty"`
	Notes                 *string                `json:"notes,omitempty"`
	Status                string                 `json:"status"`
	CurrentStepInstanceID *string                `json:"current_step_instance_id,omitempty"`
	SubmittedAt           *time.Time             `json:"submitted_at,omitempty"`
	CompletedAt           *time.Time             `json:"completed_at,omitempty"`
	CancelledAt           *time.Time             `json:"cancelled_at,omitempty"`
	CreatedAt             time.Time              `json:"created_at"`
	UpdatedAt             time.Time              `json:"updated_at"`
	ApprovalSteps         []ApprovalStepInstance `json:"approval_steps,omitempty"`
	ApprovalActions       []ApprovalAction       `json:"approval_actions,omitempty"`
}

type ApprovalStepInstance struct {
	ID                string     `json:"id"`
	OrganizationID    string     `json:"organization_id"`
	PurchaseRequestID string     `json:"purchase_request_id"`
	WorkflowVersionID string     `json:"workflow_version_id"`
	StepKey           string     `json:"step_key"`
	StepName          string     `json:"step_name"`
	SequenceNumber    int        `json:"sequence_number"`
	ApproverRoleKey   string     `json:"approver_role_key"`
	Scope             string     `json:"scope"`
	Status            string     `json:"status"`
	StartedAt         *time.Time `json:"started_at,omitempty"`
	CompletedAt       *time.Time `json:"completed_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type ApprovalAction struct {
	ID                     string    `json:"id"`
	OrganizationID         string    `json:"organization_id"`
	PurchaseRequestID      string    `json:"purchase_request_id"`
	ApprovalStepInstanceID string    `json:"approval_step_instance_id"`
	ActorUserID            string    `json:"actor_user_id"`
	Action                 string    `json:"action"`
	Comment                *string   `json:"comment,omitempty"`
	CreatedAt              time.Time `json:"created_at"`
}

type ListFilters struct {
	Page         int
	PageSize     int
	Status       string
	DepartmentID string
	RequesterID  string
	FromDate     string
	ToDate       string
	Search       string
}

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

type ListResult struct {
	Items      []PurchaseRequest `json:"items"`
	Pagination Pagination        `json:"pagination"`
}

type CreateRequest struct {
	DepartmentID    string  `json:"department_id"`
	Title           string  `json:"title"`
	ItemName        string  `json:"item_name"`
	Description     string  `json:"description"`
	Quantity        float64 `json:"quantity"`
	EstimatedAmount float64 `json:"estimated_amount"`
	Currency        string  `json:"currency"`
	Urgency         string  `json:"urgency"`
	ExpectedDate    string  `json:"expected_date"`
	VendorName      string  `json:"vendor_name"`
	Notes           string  `json:"notes"`
}

type UpdateRequest struct {
	DepartmentID    *string  `json:"department_id"`
	Title           *string  `json:"title"`
	ItemName        *string  `json:"item_name"`
	Description     *string  `json:"description"`
	Quantity        *float64 `json:"quantity"`
	EstimatedAmount *float64 `json:"estimated_amount"`
	Currency        *string  `json:"currency"`
	Urgency         *string  `json:"urgency"`
	ExpectedDate    *string  `json:"expected_date"`
	VendorName      *string  `json:"vendor_name"`
	Notes           *string  `json:"notes"`
}

type SubmitResult struct {
	PurchaseRequestID string `json:"purchase_request_id"`
	WorkflowID        string `json:"workflow_id"`
	WorkflowVersionID string `json:"workflow_version_id"`
	ApprovalSteps     int    `json:"approval_steps"`
	Status            string `json:"status"`
}

type requestRowScanner interface {
	Scan(dest ...any) error
}

func (s *Service) List(ctx context.Context, organizationID string, canViewAll bool, userID string, filters ListFilters) (*ListResult, error) {
	page, pageSize := normalizePage(filters.Page, filters.PageSize)

	where := []string{"pr.organization_id = $1"}
	args := []any{organizationID}
	argPos := 2

	if !canViewAll {
		where = append(where, fmt.Sprintf("pr.requester_id = $%d", argPos))
		args = append(args, userID)
		argPos++
	}
	if filters.Status != "" {
		where = append(where, fmt.Sprintf("pr.status = $%d", argPos))
		args = append(args, filters.Status)
		argPos++
	}
	if filters.DepartmentID != "" {
		where = append(where, fmt.Sprintf("pr.department_id = $%d", argPos))
		args = append(args, filters.DepartmentID)
		argPos++
	}
	if filters.RequesterID != "" {
		where = append(where, fmt.Sprintf("pr.requester_id = $%d", argPos))
		args = append(args, filters.RequesterID)
		argPos++
	}
	if filters.FromDate != "" {
		where = append(where, fmt.Sprintf("pr.created_at >= $%d::timestamptz", argPos))
		args = append(args, filters.FromDate)
		argPos++
	}
	if filters.ToDate != "" {
		where = append(where, fmt.Sprintf("pr.created_at <= $%d::timestamptz", argPos))
		args = append(args, filters.ToDate)
		argPos++
	}
	if filters.Search != "" {
		where = append(where, fmt.Sprintf("(pr.title ILIKE $%d OR pr.item_name ILIKE $%d OR pr.description ILIKE $%d)", argPos, argPos+1, argPos+2))
		term := "%" + filters.Search + "%"
		args = append(args, term, term, term)
		argPos += 3
	}

	limitPos := argPos
	offsetPos := argPos + 1
	args = append(args, pageSize, (page-1)*pageSize)

	query := fmt.Sprintf(`
		SELECT
			pr.id::text,
			pr.organization_id::text,
			pr.workflow_id::text,
			pr.workflow_version_id::text,
			pr.requester_id::text,
			pr.department_id::text,
			pr.title,
			pr.item_name,
			pr.description,
			pr.quantity,
			pr.estimated_amount,
			pr.currency,
			pr.urgency,
			pr.expected_date,
			pr.vendor_name,
			pr.notes,
			pr.status,
			pr.current_step_instance_id::text,
			pr.submitted_at,
			pr.completed_at,
			pr.cancelled_at,
			pr.created_at,
			pr.updated_at
		FROM purchase_requests pr
		WHERE %s
		ORDER BY pr.created_at DESC
		LIMIT $%d OFFSET $%d
	`, strings.Join(where, " AND "), limitPos, offsetPos)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]PurchaseRequest, 0)
	for rows.Next() {
		item, err := scanPurchaseRequest(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var total int
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM purchase_requests pr WHERE %s`, strings.Join(where, " AND "))
	if err := s.db.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total); err != nil {
		return nil, err
	}

	return &ListResult{
		Items: items,
		Pagination: Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}, nil
}

func (s *Service) Get(ctx context.Context, organizationID, requestID string) (*PurchaseRequest, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT
			pr.id::text,
			pr.organization_id::text,
			pr.workflow_id::text,
			pr.workflow_version_id::text,
			pr.requester_id::text,
			pr.department_id::text,
			pr.title,
			pr.item_name,
			pr.description,
			pr.quantity,
			pr.estimated_amount,
			pr.currency,
			pr.urgency,
			pr.expected_date,
			pr.vendor_name,
			pr.notes,
			pr.status,
			pr.current_step_instance_id::text,
			pr.submitted_at,
			pr.completed_at,
			pr.cancelled_at,
			pr.created_at,
			pr.updated_at
		FROM purchase_requests pr
		WHERE pr.organization_id = $1 AND pr.id = $2
	`, organizationID, requestID)

	item, err := scanPurchaseRequest(row)
	if err != nil {
		return nil, err
	}

	steps, err := s.listSteps(ctx, organizationID, requestID)
	if err != nil {
		return nil, err
	}
	actions, err := s.listActions(ctx, organizationID, requestID)
	if err != nil {
		return nil, err
	}

	item.ApprovalSteps = steps
	item.ApprovalActions = actions
	return item, nil
}

func (s *Service) Create(ctx context.Context, organizationID, requesterID string, req CreateRequest) (*PurchaseRequest, error) {
	if err := validateCreateRequest(req); err != nil {
		return nil, err
	}

	departmentID := strings.TrimSpace(req.DepartmentID)
	if err := s.ensureDepartment(ctx, organizationID, departmentID); err != nil {
		return nil, err
	}

	expectedDate, err := parseDate(req.ExpectedDate)
	if err != nil {
		return nil, err
	}

	row := s.db.QueryRowContext(ctx, `
		INSERT INTO purchase_requests (
			organization_id, requester_id, department_id, title, item_name, description,
			quantity, estimated_amount, currency, urgency, expected_date, vendor_name, notes, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, COALESCE(NULLIF($9, ''), 'IDR'), COALESCE(NULLIF($10, ''), 'normal'), $11, NULLIF($12, ''), NULLIF($13, ''), 'draft')
		RETURNING
			id::text,
			organization_id::text,
			workflow_id::text,
			workflow_version_id::text,
			requester_id::text,
			department_id::text,
			title,
			item_name,
			description,
			quantity,
			estimated_amount,
			currency,
			urgency,
			expected_date,
			vendor_name,
			notes,
			status,
			current_step_instance_id::text,
			submitted_at,
			completed_at,
			cancelled_at,
			created_at,
			updated_at
	`,
		organizationID,
		requesterID,
		departmentID,
		strings.TrimSpace(req.Title),
		strings.TrimSpace(req.ItemName),
		strings.TrimSpace(req.Description),
		req.Quantity,
		req.EstimatedAmount,
		strings.TrimSpace(req.Currency),
		strings.TrimSpace(req.Urgency),
		dateArg(expectedDate),
		strings.TrimSpace(req.VendorName),
		strings.TrimSpace(req.Notes),
	)
	if err != nil {
		return nil, err
	}

	return scanPurchaseRequest(row)
}

func (s *Service) Update(ctx context.Context, organizationID, requesterID, requestID string, req UpdateRequest) (*PurchaseRequest, error) {
	current, err := s.Get(ctx, organizationID, requestID)
	if err != nil {
		return nil, err
	}

	if current.RequesterID != requesterID && requesterID != "" {
		// allow privileged callers to update via organization scoping in handlers later
	}

	if current.Status != "draft" && current.Status != "revision_requested" {
		return nil, errors.New("request can only be updated in draft or revision_requested status")
	}

	departmentID := current.DepartmentID
	if req.DepartmentID != nil {
		departmentID = strings.TrimSpace(*req.DepartmentID)
		if err := s.ensureDepartment(ctx, organizationID, departmentID); err != nil {
			return nil, err
		}
	}
	if err := s.ensureDepartment(ctx, organizationID, departmentID); err != nil {
		return nil, err
	}

	title := current.Title
	if req.Title != nil {
		title = strings.TrimSpace(*req.Title)
	}
	itemName := current.ItemName
	if req.ItemName != nil {
		itemName = strings.TrimSpace(*req.ItemName)
	}
	description := current.Description
	if req.Description != nil {
		description = strings.TrimSpace(*req.Description)
	}
	quantity := current.Quantity
	if req.Quantity != nil {
		quantity = *req.Quantity
	}
	estimatedAmount := current.EstimatedAmount
	if req.EstimatedAmount != nil {
		estimatedAmount = *req.EstimatedAmount
	}
	currency := current.Currency
	if req.Currency != nil {
		currency = strings.TrimSpace(*req.Currency)
	}
	urgency := current.Urgency
	if req.Urgency != nil {
		urgency = strings.TrimSpace(*req.Urgency)
	}
	var expectedDateArg *time.Time
	if current.ExpectedDate != nil {
		parsed, err := time.Parse("2006-01-02", *current.ExpectedDate)
		if err != nil {
			return nil, err
		}
		expectedDateArg = &parsed
	}
	if req.ExpectedDate != nil {
		parsed, err := parseDate(*req.ExpectedDate)
		if err != nil {
			return nil, err
		}
		expectedDateArg = parsed
	}
	vendorName := current.VendorName
	if req.VendorName != nil {
		trimmed := strings.TrimSpace(*req.VendorName)
		if trimmed == "" {
			vendorName = nil
		} else {
			vendorName = &trimmed
		}
	}
	notes := current.Notes
	if req.Notes != nil {
		trimmed := strings.TrimSpace(*req.Notes)
		if trimmed == "" {
			notes = nil
		} else {
			notes = &trimmed
		}
	}

	row := s.db.QueryRowContext(ctx, `
		UPDATE purchase_requests
		SET department_id = $3,
		    title = $4,
		    item_name = $5,
		    description = $6,
		    quantity = $7,
		    estimated_amount = $8,
		    currency = $9,
		    urgency = $10,
		    expected_date = $11,
		    vendor_name = $12,
		    notes = $13,
		    updated_at = NOW()
		WHERE organization_id = $1 AND id = $2
		RETURNING
			id::text,
			organization_id::text,
			workflow_id::text,
			workflow_version_id::text,
			requester_id::text,
			department_id::text,
			title,
			item_name,
			description,
			quantity,
			estimated_amount,
			currency,
			urgency,
			expected_date,
			vendor_name,
			notes,
			status,
			current_step_instance_id::text,
			submitted_at,
			completed_at,
			cancelled_at,
			created_at,
			updated_at
		`,
		organizationID,
		requestID,
		departmentID,
		title,
		itemName,
		description,
		quantity,
		estimatedAmount,
		currency,
		urgency,
		dateArg(expectedDateArg),
		nullableString(vendorName),
		nullableString(notes),
	)
	if err != nil {
		return nil, err
	}

	return scanPurchaseRequest(row)
}

func (s *Service) Submit(ctx context.Context, organizationID, requesterID, requestID string) (*SubmitResult, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	current, err := s.getForUpdate(ctx, tx, organizationID, requestID)
	if err != nil {
		return nil, err
	}

	if current.RequesterID != requesterID {
		return nil, errors.New("only requester can submit request")
	}
	if current.Status != "draft" && current.Status != "revision_requested" {
		return nil, errors.New("request can only be submitted from draft or revision_requested")
	}

	workflow, err := s.getActiveWorkflow(ctx, tx, organizationID)
	if err != nil {
		return nil, err
	}

	steps, err := s.generateSteps(ctx, tx, current, workflow)
	if err != nil {
		return nil, err
	}
	if len(steps) == 0 {
		return nil, errors.New("workflow produced no applicable steps")
	}

	firstStepID := steps[0].ID
	if _, err := tx.ExecContext(ctx, `
		UPDATE purchase_requests
		SET workflow_id = $3,
		    workflow_version_id = $4,
		    current_step_instance_id = $5,
		    status = 'in_review',
		    submitted_at = NOW(),
		    updated_at = NOW()
		WHERE organization_id = $1 AND id = $2
	`, organizationID, requestID, workflow.WorkflowID, workflow.VersionID, firstStepID); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &SubmitResult{
		PurchaseRequestID: requestID,
		WorkflowID:        workflow.WorkflowID,
		WorkflowVersionID: workflow.VersionID,
		ApprovalSteps:     len(steps),
		Status:            "in_review",
	}, nil
}

func (s *Service) getForUpdate(ctx context.Context, tx *sql.Tx, organizationID, requestID string) (*PurchaseRequest, error) {
	row := tx.QueryRowContext(ctx, `
		SELECT
			id::text,
			organization_id::text,
			workflow_id::text,
			workflow_version_id::text,
			requester_id::text,
			department_id::text,
			title,
			item_name,
			description,
			quantity,
			estimated_amount,
			currency,
			urgency,
			expected_date,
			vendor_name,
			notes,
			status,
			current_step_instance_id::text,
			submitted_at,
			completed_at,
			cancelled_at,
			created_at,
			updated_at
		FROM purchase_requests
		WHERE organization_id = $1 AND id = $2
		FOR UPDATE
	`, organizationID, requestID)

	return scanPurchaseRequest(row)
}

func (s *Service) getActiveWorkflow(ctx context.Context, tx *sql.Tx, organizationID string) (*WorkflowLookup, error) {
	row := tx.QueryRowContext(ctx, `
		SELECT
			w.id::text,
			w.active_version_id::text,
			v.config_json
		FROM workflows w
		INNER JOIN workflow_versions v ON v.id = w.active_version_id
		WHERE w.organization_id = $1
		  AND w.type = 'purchase_request'
		  AND w.status = 'active'
		ORDER BY w.updated_at DESC
		LIMIT 1
	`, organizationID)

	var workflowID, versionID string
	var cfgBytes []byte
	if err := row.Scan(&workflowID, &versionID, &cfgBytes); err != nil {
		return nil, err
	}

	var cfg configworkflow.Config
	if err := json.Unmarshal(cfgBytes, &cfg); err != nil {
		return nil, err
	}
	if validation := configworkflow.ValidateConfig(cfg); !validation.Valid {
		return nil, ErrInvalidWorkflow{Validation: validation}
	}

	return &WorkflowLookup{WorkflowID: workflowID, VersionID: versionID, Config: cfg}, nil
}

func (s *Service) generateSteps(ctx context.Context, tx *sql.Tx, request *PurchaseRequest, workflow *WorkflowLookup) ([]ApprovalStepInstance, error) {
	applicable, err := applicableSteps(workflow.Config, request)
	if err != nil {
		return nil, err
	}

	instances := make([]ApprovalStepInstance, 0, len(applicable))
	for i, step := range applicable {
		status := "waiting"
		if i == 0 {
			status = "pending"
		}
		instance := ApprovalStepInstance{}
		err := tx.QueryRowContext(ctx, `
			INSERT INTO approval_step_instances (
				organization_id, purchase_request_id, workflow_version_id, step_key, step_name,
				sequence_number, approver_role_key, scope, status, started_at, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, CASE WHEN $9 = 'pending' THEN NOW() ELSE NULL END, NOW(), NOW())
			RETURNING
				id::text,
				organization_id::text,
				purchase_request_id::text,
				workflow_version_id::text,
				step_key,
				step_name,
				sequence_number,
				approver_role_key,
				scope,
				status,
				started_at,
				completed_at,
				created_at,
				updated_at
		`, request.OrganizationID, request.ID, workflow.VersionID, step.Key, step.Name, i+1, step.ApproverRole, step.Scope, status).Scan(
			&instance.ID,
			&instance.OrganizationID,
			&instance.PurchaseRequestID,
			&instance.WorkflowVersionID,
			&instance.StepKey,
			&instance.StepName,
			&instance.SequenceNumber,
			&instance.ApproverRoleKey,
			&instance.Scope,
			&instance.Status,
			&instance.StartedAt,
			&instance.CompletedAt,
			&instance.CreatedAt,
			&instance.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		instances = append(instances, instance)
	}

	return instances, nil
}

func applicableSteps(cfg configworkflow.Config, request *PurchaseRequest) ([]configworkflow.Step, error) {
	applicable := make([]configworkflow.Step, 0, len(cfg.Steps))
	for _, step := range cfg.Steps {
		if step.Condition == nil {
			applicable = append(applicable, step)
			continue
		}

		matches, err := evaluateStepCondition(step, request)
		if err != nil {
			return nil, err
		}
		if matches {
			applicable = append(applicable, step)
		}
	}
	return applicable, nil
}

func (s *Service) listSteps(ctx context.Context, organizationID, requestID string) ([]ApprovalStepInstance, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT
			id::text,
			organization_id::text,
			purchase_request_id::text,
			workflow_version_id::text,
			step_key,
			step_name,
			sequence_number,
			approver_role_key,
			scope,
			status,
			started_at,
			completed_at,
			created_at,
			updated_at
		FROM approval_step_instances
		WHERE organization_id = $1 AND purchase_request_id = $2
		ORDER BY sequence_number ASC
	`, organizationID, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []ApprovalStepInstance
	for rows.Next() {
		var item ApprovalStepInstance
		if err := rows.Scan(
			&item.ID,
			&item.OrganizationID,
			&item.PurchaseRequestID,
			&item.WorkflowVersionID,
			&item.StepKey,
			&item.StepName,
			&item.SequenceNumber,
			&item.ApproverRoleKey,
			&item.Scope,
			&item.Status,
			&item.StartedAt,
			&item.CompletedAt,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Service) listActions(ctx context.Context, organizationID, requestID string) ([]ApprovalAction, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT
			id::text,
			organization_id::text,
			purchase_request_id::text,
			approval_step_instance_id::text,
			actor_user_id::text,
			action,
			comment,
			created_at
		FROM approval_actions
		WHERE organization_id = $1 AND purchase_request_id = $2
		ORDER BY created_at ASC
	`, organizationID, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []ApprovalAction
	for rows.Next() {
		var item ApprovalAction
		var comment sql.NullString
		if err := rows.Scan(
			&item.ID,
			&item.OrganizationID,
			&item.PurchaseRequestID,
			&item.ApprovalStepInstanceID,
			&item.ActorUserID,
			&item.Action,
			&comment,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}
		if comment.Valid {
			value := comment.String
			item.Comment = &value
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Service) ensureDepartment(ctx context.Context, organizationID, departmentID string) error {
	if strings.TrimSpace(departmentID) == "" {
		return errors.New("department_id is required")
	}
	var id string
	if err := s.db.QueryRowContext(ctx, `
		SELECT id::text
		FROM departments
		WHERE organization_id = $1 AND id = $2
	`, organizationID, departmentID).Scan(&id); err != nil {
		return err
	}
	return nil
}

func validateCreateRequest(req CreateRequest) error {
	switch {
	case strings.TrimSpace(req.DepartmentID) == "":
		return errors.New("department_id is required")
	case strings.TrimSpace(req.Title) == "":
		return errors.New("title is required")
	case strings.TrimSpace(req.ItemName) == "":
		return errors.New("item_name is required")
	case strings.TrimSpace(req.Description) == "":
		return errors.New("description is required")
	case req.Quantity <= 0:
		return errors.New("quantity must be greater than 0")
	case req.EstimatedAmount <= 0:
		return errors.New("estimated_amount must be greater than 0")
	}
	return nil
}

func parseDate(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, errors.New("expected_date must be YYYY-MM-DD")
	}
	return &parsed, nil
}

func nullableString(value *string) any {
	if value == nil {
		return nil
	}
	return *value
}

func dateArg(value *time.Time) any {
	if value == nil {
		return nil
	}
	return *value
}

func normalizePage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func scanPurchaseRequest(scanner requestRowScanner) (*PurchaseRequest, error) {
	var item PurchaseRequest
	var workflowID, workflowVersionID, currentStepID sql.NullString
	var departmentID string
	var expectedDate sql.NullTime
	var vendorName, notes sql.NullString
	var submittedAt, completedAt, cancelledAt sql.NullTime
	if err := scanner.Scan(
		&item.ID,
		&item.OrganizationID,
		&workflowID,
		&workflowVersionID,
		&item.RequesterID,
		&departmentID,
		&item.Title,
		&item.ItemName,
		&item.Description,
		&item.Quantity,
		&item.EstimatedAmount,
		&item.Currency,
		&item.Urgency,
		&expectedDate,
		&vendorName,
		&notes,
		&item.Status,
		&currentStepID,
		&submittedAt,
		&completedAt,
		&cancelledAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return nil, err
	}

	if workflowID.Valid {
		value := workflowID.String
		item.WorkflowID = &value
	}
	if workflowVersionID.Valid {
		value := workflowVersionID.String
		item.WorkflowVersionID = &value
	}
	if currentStepID.Valid {
		value := currentStepID.String
		item.CurrentStepInstanceID = &value
	}
	if expectedDate.Valid {
		value := expectedDate.Time.Format("2006-01-02")
		item.ExpectedDate = &value
	}
	if vendorName.Valid {
		value := vendorName.String
		item.VendorName = &value
	}
	if notes.Valid {
		value := notes.String
		item.Notes = &value
	}
	if submittedAt.Valid {
		value := submittedAt.Time
		item.SubmittedAt = &value
	}
	if completedAt.Valid {
		value := completedAt.Time
		item.CompletedAt = &value
	}
	if cancelledAt.Valid {
		value := cancelledAt.Time
		item.CancelledAt = &value
	}

	return &item, nil
}

func evaluateStepCondition(step configworkflow.Step, request *PurchaseRequest) (bool, error) {
	if step.Condition == nil {
		return true, nil
	}

	var value any
	switch step.Condition.Field {
	case "estimated_amount":
		value = request.EstimatedAmount
	case "currency":
		value = request.Currency
	case "urgency":
		value = request.Urgency
	case "department_id":
		value = request.DepartmentID
	default:
		return false, errors.New("condition field is not supported")
	}

	return configworkflow.EvaluateCondition(*step.Condition, value)
}

type ErrInvalidWorkflow struct {
	Validation configworkflow.ValidationResult
}

func (e ErrInvalidWorkflow) Error() string {
	return "workflow invalid"
}

var ErrConflict = errors.New("conflict")
