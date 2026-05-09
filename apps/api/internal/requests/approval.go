package requests

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

type ActorContext struct {
	UserID       string
	DepartmentID *string
	Roles        []string
}

type ApprovalActionRequest struct {
	Action  string `json:"action"`
	Comment string `json:"comment"`
}

type ApprovalActionResult struct {
	PurchaseRequestID      string     `json:"purchase_request_id"`
	ApprovalStepInstanceID string     `json:"approval_step_instance_id"`
	Action                 string     `json:"action"`
	Status                 string     `json:"status"`
	CurrentStepInstanceID  *string    `json:"current_step_instance_id,omitempty"`
	CompletedAt            *time.Time `json:"completed_at,omitempty"`
}

type pendingApprovalRow struct {
	RequestID string `json:"request_id"`
	StepID    string `json:"step_id"`
	StepName  string `json:"step_name"`
	RoleKey   string `json:"role_key"`
	Scope     string `json:"scope"`
	Status    string `json:"status"`
}

func (s *Service) ListPendingApprovals(ctx context.Context, organizationID string, actor ActorContext) ([]pendingApprovalRow, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT
			pr.id::text,
			asi.id::text,
			asi.step_name,
			asi.approver_role_key,
			asi.scope,
			asi.status
		FROM purchase_requests pr
		INNER JOIN approval_step_instances asi ON asi.id = pr.current_step_instance_id
		WHERE pr.organization_id = $1
		  AND pr.status = 'in_review'
		  AND asi.status = 'pending'
		  AND (
			asi.scope <> 'requester_department'
			OR EXISTS (
				SELECT 1
				FROM users u
				WHERE u.id = $2
				  AND u.organization_id = $1
				  AND u.department_id = pr.department_id
			)
		  )
		  AND EXISTS (
			SELECT 1
			FROM user_roles ur
			INNER JOIN roles r ON r.id = ur.role_id
			WHERE ur.organization_id = $1
			  AND ur.user_id = $2
			  AND r.key = asi.approver_role_key
		  )
		ORDER BY pr.created_at DESC
	`, organizationID, actor.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []pendingApprovalRow
	for rows.Next() {
		var item pendingApprovalRow
		if err := rows.Scan(&item.RequestID, &item.StepID, &item.StepName, &item.RoleKey, &item.Scope, &item.Status); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Service) ActOnRequest(ctx context.Context, organizationID string, actor ActorContext, requestID string, req ApprovalActionRequest) (*ApprovalActionResult, error) {
	action := strings.TrimSpace(strings.ToLower(req.Action))
	if action == "" {
		return nil, errors.New("action is required")
	}
	if action != "approve" && action != "reject" && action != "request_revision" {
		return nil, errors.New("action is not supported")
	}
	if action == "reject" || action == "request_revision" {
		if strings.TrimSpace(req.Comment) == "" {
			return nil, errors.New("comment is required")
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	current, err := s.getForUpdate(ctx, tx, organizationID, requestID)
	if err != nil {
		return nil, err
	}
	if current.Status != "in_review" && current.Status != "processing" {
		return nil, ErrForbidden
	}
	if current.CurrentStepInstanceID == nil {
		return nil, ErrForbidden
	}
	if current.RequesterID == actor.UserID {
		return nil, ErrForbidden
	}

	step, err := s.getStepForUpdate(ctx, tx, organizationID, requestID, *current.CurrentStepInstanceID)
	if err != nil {
		return nil, err
	}
	if step.Status != "pending" {
		return nil, ErrForbidden
	}
	if !hasRole(actor.Roles, step.ApproverRoleKey) {
		return nil, ErrForbidden
	}
	if step.Scope == "requester_department" {
		if actor.DepartmentID == nil || strings.TrimSpace(*actor.DepartmentID) == "" {
			return nil, ErrForbidden
		}
		if strings.TrimSpace(*actor.DepartmentID) != strings.TrimSpace(current.DepartmentID) {
			return nil, ErrForbidden
		}
	}

	var nextStep *ApprovalStepInstance
	if action == "approve" {
		nextStep, err = s.nextStepForUpdate(ctx, tx, organizationID, requestID, step.SequenceNumber)
		if err != nil {
			return nil, err
		}
	}

	now := time.Now().UTC()
	actionRecord := map[string]any{
		"action":                 action,
		"comment":                strings.TrimSpace(req.Comment),
		"approval_step_instance": step.ID,
		"purchase_request_id":    current.ID,
	}

	if err := s.recordAuditTx(ctx, tx, organizationID, actor.UserID, "approval."+action, "approval_step_instance", &step.ID, nil, actionRecord); err != nil {
		return nil, err
	}

	var status string
	switch action {
	case "approve":
		status = "approved"
	case "reject":
		status = "rejected"
	case "request_revision":
		status = "revision_requested"
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE approval_step_instances
		SET status = $4,
		    completed_at = $5,
		    updated_at = NOW()
		WHERE organization_id = $1 AND purchase_request_id = $2 AND id = $3
	`, organizationID, requestID, step.ID, status, now); err != nil {
		return nil, err
	}

	var currentStepID any
	var requestStatus string
	var completedAt any

	switch action {
	case "approve":
		if nextStep != nil {
			if _, err := tx.ExecContext(ctx, `
				UPDATE approval_step_instances
				SET status = 'pending',
				    started_at = COALESCE(started_at, NOW()),
				    updated_at = NOW()
				WHERE organization_id = $1 AND purchase_request_id = $2 AND id = $3
			`, organizationID, requestID, nextStep.ID); err != nil {
				return nil, err
			}
			currentStepID = nextStep.ID
			requestStatus = "in_review"
		} else {
			currentStepID = nil
			if strings.Contains(strings.ToLower(step.ApproverRoleKey), "procurement") || strings.Contains(strings.ToLower(step.StepKey), "procurement") {
				requestStatus = "completed"
				completedAt = now
			} else {
				requestStatus = "approved"
			}
		}
	case "reject":
		currentStepID = nil
		requestStatus = "rejected"
		completedAt = now
	case "request_revision":
		currentStepID = nil
		requestStatus = "revision_requested"
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE purchase_requests
		SET status = $3,
		    current_step_instance_id = $4,
		    completed_at = COALESCE($5, completed_at),
		    updated_at = NOW()
		WHERE organization_id = $1 AND id = $2
	`, organizationID, requestID, requestStatus, currentStepID, completedAt); err != nil {
		return nil, err
	}

	if err := s.recordAuditTx(ctx, tx, organizationID, actor.UserID, "request."+status, "purchase_request", &requestID, current, map[string]any{
		"status":                    requestStatus,
		"current_step_instance_id":  currentStepID,
		"approval_step_instance_id": step.ID,
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	var currentStepPtr *string
	if v, ok := currentStepID.(string); ok && v != "" {
		currentStepPtr = &v
	}
	var completedPtr *time.Time
	if t, ok := completedAt.(time.Time); ok {
		completedPtr = &t
	}

	return &ApprovalActionResult{
		PurchaseRequestID:      requestID,
		ApprovalStepInstanceID: step.ID,
		Action:                 status,
		Status:                 requestStatus,
		CurrentStepInstanceID:  currentStepPtr,
		CompletedAt:            completedPtr,
	}, nil
}

func (s *Service) AddComment(ctx context.Context, organizationID string, actor ActorContext, requestID, body string) (*Comment, error) {
	body = strings.TrimSpace(body)
	if body == "" {
		return nil, errors.New("body is required")
	}

	current, err := s.Get(ctx, organizationID, requestID)
	if err != nil {
		return nil, err
	}
	if current.RequesterID != actor.UserID && !hasAnyRole(actor.Roles, "owner", "admin", "auditor", "finance", "manager", "director", "procurement") {
		return nil, ErrForbidden
	}

	row := s.db.QueryRowContext(ctx, `
		INSERT INTO comments (organization_id, purchase_request_id, actor_user_id, body)
		VALUES ($1, $2, $3, $4)
		RETURNING id::text, organization_id::text, purchase_request_id::text, actor_user_id::text, body, created_at, updated_at
	`, organizationID, requestID, actor.UserID, body)

	var comment Comment
	if err := row.Scan(&comment.ID, &comment.OrganizationID, &comment.PurchaseRequestID, &comment.ActorUserID, &comment.Body, &comment.CreatedAt, &comment.UpdatedAt); err != nil {
		return nil, err
	}

	if err := s.recordAudit(ctx, organizationID, actor.UserID, "request.commented", "purchase_request", &requestID, nil, comment); err != nil {
		return nil, err
	}

	return &comment, nil
}

func (s *Service) ListComments(ctx context.Context, organizationID, requestID string) ([]Comment, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id::text, organization_id::text, purchase_request_id::text, actor_user_id::text, body, created_at, updated_at
		FROM comments
		WHERE organization_id = $1 AND purchase_request_id = $2
		ORDER BY created_at ASC
	`, organizationID, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Comment
	for rows.Next() {
		var item Comment
		if err := rows.Scan(&item.ID, &item.OrganizationID, &item.PurchaseRequestID, &item.ActorUserID, &item.Body, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

type Comment struct {
	ID                string    `json:"id"`
	OrganizationID    string    `json:"organization_id"`
	PurchaseRequestID string    `json:"purchase_request_id"`
	ActorUserID       string    `json:"actor_user_id"`
	Body              string    `json:"body"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (s *Service) getStepForUpdate(ctx context.Context, tx *sql.Tx, organizationID, requestID, stepID string) (*ApprovalStepInstance, error) {
	row := tx.QueryRowContext(ctx, `
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
		WHERE organization_id = $1 AND purchase_request_id = $2 AND id = $3
		FOR UPDATE
	`, organizationID, requestID, stepID)

	var step ApprovalStepInstance
	if err := row.Scan(
		&step.ID,
		&step.OrganizationID,
		&step.PurchaseRequestID,
		&step.WorkflowVersionID,
		&step.StepKey,
		&step.StepName,
		&step.SequenceNumber,
		&step.ApproverRoleKey,
		&step.Scope,
		&step.Status,
		&step.StartedAt,
		&step.CompletedAt,
		&step.CreatedAt,
		&step.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &step, nil
}

func (s *Service) nextStepForUpdate(ctx context.Context, tx *sql.Tx, organizationID, requestID string, sequenceNumber int) (*ApprovalStepInstance, error) {
	row := tx.QueryRowContext(ctx, `
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
		WHERE organization_id = $1
		  AND purchase_request_id = $2
		  AND sequence_number = $3 + 1
		FOR UPDATE
	`, organizationID, requestID, sequenceNumber)

	var step ApprovalStepInstance
	if err := row.Scan(
		&step.ID,
		&step.OrganizationID,
		&step.PurchaseRequestID,
		&step.WorkflowVersionID,
		&step.StepKey,
		&step.StepName,
		&step.SequenceNumber,
		&step.ApproverRoleKey,
		&step.Scope,
		&step.Status,
		&step.StartedAt,
		&step.CompletedAt,
		&step.CreatedAt,
		&step.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &step, nil
}

func hasRole(roles []string, target string) bool {
	return hasAnyRole(roles, target)
}

func hasAnyRole(roles []string, targets ...string) bool {
	roleSet := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		roleSet[role] = struct{}{}
	}
	for _, target := range targets {
		if _, ok := roleSet[target]; ok {
			return true
		}
	}
	return false
}
