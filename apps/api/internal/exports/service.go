package exports

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"operra/api/internal/audit"
)

type Service struct {
	db    *sql.DB
	audit *audit.Service
}

func NewService(db *sql.DB, auditSvc ...*audit.Service) *Service {
	var svc *audit.Service
	if len(auditSvc) > 0 {
		svc = auditSvc[0]
	}
	return &Service{db: db, audit: svc}
}

type ExportFilters struct {
	Status       string
	DepartmentID string
	FromDate     string
	ToDate       string
	RequestID    string
	Action       string
	EntityType   string
	ActorUserID  string
}

func (s *Service) PurchaseRequestsCSV(ctx context.Context, organizationID string, canViewAll bool, userID string, filters ExportFilters) (string, int, error) {
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

	query := fmt.Sprintf(`
		SELECT
			pr.id::text,
			pr.title,
			d.name,
			u.name,
			pr.estimated_amount,
			pr.currency,
			pr.status,
			pr.submitted_at,
			pr.completed_at,
			pr.created_at
		FROM purchase_requests pr
		INNER JOIN users u ON u.id = pr.requester_id AND u.organization_id = pr.organization_id
		INNER JOIN departments d ON d.id = pr.department_id AND d.organization_id = pr.organization_id
		WHERE %s
		ORDER BY pr.created_at DESC
	`, strings.Join(where, " AND "))

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return "", 0, err
	}
	defer rows.Close()

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	if err := writer.Write([]string{"request_id", "title", "department", "requester_name", "estimated_amount", "currency", "status", "submitted_at", "completed_at", "created_at"}); err != nil {
		return "", 0, err
	}

	count := 0
	for rows.Next() {
		var requestID, title, department, requester, currency, status string
		var estimatedAmount float64
		var submittedAt, completedAt, createdAt sql.NullTime
		if err := rows.Scan(&requestID, &title, &department, &requester, &estimatedAmount, &currency, &status, &submittedAt, &completedAt, &createdAt); err != nil {
			return "", 0, err
		}
		record := []string{
			requestID,
			sanitizeCSV(title),
			sanitizeCSV(department),
			sanitizeCSV(requester),
			strconv.FormatFloat(estimatedAmount, 'f', 2, 64),
			currency,
			status,
			formatTime(submittedAt),
			formatTime(completedAt),
			formatTime(createdAt),
		}
		if err := writer.Write(record); err != nil {
			return "", 0, err
		}
		count++
	}
	if err := rows.Err(); err != nil {
		return "", 0, err
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", 0, err
	}

	if err := s.recordExport(ctx, organizationID, userID, "purchase_requests", filters, count); err != nil {
		return "", 0, err
	}

	return buf.String(), count, nil
}

func (s *Service) ApprovalHistoryCSV(ctx context.Context, organizationID string, userID string, filters ExportFilters) (string, int, error) {
	where := []string{"a.organization_id = $1"}
	args := []any{organizationID}
	argPos := 2

	if filters.RequestID != "" {
		where = append(where, fmt.Sprintf("a.purchase_request_id = $%d", argPos))
		args = append(args, filters.RequestID)
		argPos++
	}
	if filters.FromDate != "" {
		where = append(where, fmt.Sprintf("a.created_at >= $%d::timestamptz", argPos))
		args = append(args, filters.FromDate)
		argPos++
	}
	if filters.ToDate != "" {
		where = append(where, fmt.Sprintf("a.created_at <= $%d::timestamptz", argPos))
		args = append(args, filters.ToDate)
		argPos++
	}

	query := fmt.Sprintf(`
		SELECT
			a.purchase_request_id::text,
			pr.title,
			asi.step_name,
			a.action,
			u.name,
			COALESCE(a.comment, ''),
			a.created_at
		FROM approval_actions a
		INNER JOIN purchase_requests pr ON pr.id = a.purchase_request_id AND pr.organization_id = a.organization_id
		INNER JOIN approval_step_instances asi ON asi.id = a.approval_step_instance_id AND asi.organization_id = a.organization_id
		INNER JOIN users u ON u.id = a.actor_user_id AND u.organization_id = a.organization_id
		WHERE %s
		ORDER BY a.created_at DESC
	`, strings.Join(where, " AND "))

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return "", 0, err
	}
	defer rows.Close()

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	if err := writer.Write([]string{"request_id", "request_title", "step_name", "action", "actor_name", "comment", "created_at"}); err != nil {
		return "", 0, err
	}

	count := 0
	for rows.Next() {
		var requestID, title, stepName, action, actorName, comment string
		var createdAt time.Time
		if err := rows.Scan(&requestID, &title, &stepName, &action, &actorName, &comment, &createdAt); err != nil {
			return "", 0, err
		}
		record := []string{requestID, sanitizeCSV(title), sanitizeCSV(stepName), action, sanitizeCSV(actorName), sanitizeCSV(comment), createdAt.UTC().Format(time.RFC3339)}
		if err := writer.Write(record); err != nil {
			return "", 0, err
		}
		count++
	}
	if err := rows.Err(); err != nil {
		return "", 0, err
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", 0, err
	}

	if err := s.recordExport(ctx, organizationID, userID, "approval_history", filters, count); err != nil {
		return "", 0, err
	}

	return buf.String(), count, nil
}

func (s *Service) AuditLogsCSV(ctx context.Context, organizationID string, userID string, filters ExportFilters) (string, int, error) {
	where := []string{"al.organization_id = $1"}
	args := []any{organizationID}
	argPos := 2

	if filters.Action != "" {
		where = append(where, fmt.Sprintf("al.action = $%d", argPos))
		args = append(args, filters.Action)
		argPos++
	}
	if filters.EntityType != "" {
		where = append(where, fmt.Sprintf("al.entity_type = $%d", argPos))
		args = append(args, filters.EntityType)
		argPos++
	}
	if filters.FromDate != "" {
		where = append(where, fmt.Sprintf("al.created_at >= $%d::timestamptz", argPos))
		args = append(args, filters.FromDate)
		argPos++
	}
	if filters.ToDate != "" {
		where = append(where, fmt.Sprintf("al.created_at <= $%d::timestamptz", argPos))
		args = append(args, filters.ToDate)
		argPos++
	}

	query := fmt.Sprintf(`
		SELECT
			al.action,
			al.entity_type,
			COALESCE(al.entity_id::text, ''),
			COALESCE(u.name, ''),
			al.created_at
		FROM audit_logs al
		LEFT JOIN users u ON u.id = al.actor_user_id AND u.organization_id = al.organization_id
		WHERE %s
		ORDER BY al.created_at DESC
	`, strings.Join(where, " AND "))

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return "", 0, err
	}
	defer rows.Close()

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	if err := writer.Write([]string{"action", "entity_type", "entity_id", "actor_name", "created_at"}); err != nil {
		return "", 0, err
	}

	count := 0
	for rows.Next() {
		var action, entityType, entityID, actorName string
		var createdAt time.Time
		if err := rows.Scan(&action, &entityType, &entityID, &actorName, &createdAt); err != nil {
			return "", 0, err
		}
		record := []string{action, entityType, entityID, sanitizeCSV(actorName), createdAt.UTC().Format(time.RFC3339)}
		if err := writer.Write(record); err != nil {
			return "", 0, err
		}
		count++
	}
	if err := rows.Err(); err != nil {
		return "", 0, err
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", 0, err
	}

	if err := s.recordExport(ctx, organizationID, userID, "audit_logs", filters, count); err != nil {
		return "", 0, err
	}

	return buf.String(), count, nil
}

func (s *Service) recordExport(ctx context.Context, organizationID, userID, exportType string, filters ExportFilters, rowCount int) error {
	filterJSON, err := json.Marshal(filters)
	if err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO export_logs (organization_id, requested_by, export_type, filters_json, row_count)
		VALUES ($1, $2, $3, $4, $5)
	`, organizationID, userID, exportType, filterJSON, rowCount); err != nil {
		return err
	}

	if s.audit != nil {
		if err := s.audit.RecordTx(ctx, tx, organizationID, &userID, "csv.exported", "export", nil, nil, map[string]any{
			"export_type": exportType,
			"row_count":   rowCount,
			"filters":     filters,
		}, nil, nil); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func sanitizeCSV(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	switch value[0] {
	case '=', '+', '-', '@':
		return "'" + value
	default:
		return value
	}
}

func formatTime(value sql.NullTime) string {
	if !value.Valid {
		return ""
	}
	return value.Time.UTC().Format(time.RFC3339)
}
