package audit

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

type AuditLog struct {
	ID             string          `json:"id"`
	OrganizationID string          `json:"organization_id"`
	ActorUserID    *string         `json:"actor_user_id,omitempty"`
	Action         string          `json:"action"`
	EntityType     string          `json:"entity_type"`
	EntityID       *string         `json:"entity_id,omitempty"`
	OldValue       json.RawMessage `json:"old_value,omitempty"`
	NewValue       json.RawMessage `json:"new_value,omitempty"`
	IPAddr         *string         `json:"ip_address,omitempty"`
	UserAgent      *string         `json:"user_agent,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
}

type ListFilters struct {
	Page        int
	PageSize    int
	Action      string
	EntityType  string
	EntityID    string
	FromDate    string
	ToDate      string
	ActorUserID string
}

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

type ListResult struct {
	Items      []AuditLog `json:"items"`
	Pagination Pagination `json:"pagination"`
}

func (s *Service) Record(ctx context.Context, organizationID string, actorUserID *string, action, entityType string, entityID *string, oldValue, newValue any, ipAddress, userAgent *string) error {
	return s.record(ctx, s.db, organizationID, actorUserID, action, entityType, entityID, oldValue, newValue, ipAddress, userAgent)
}

func (s *Service) RecordTx(ctx context.Context, tx *sql.Tx, organizationID string, actorUserID *string, action, entityType string, entityID *string, oldValue, newValue any, ipAddress, userAgent *string) error {
	return s.record(ctx, tx, organizationID, actorUserID, action, entityType, entityID, oldValue, newValue, ipAddress, userAgent)
}

func (s *Service) List(ctx context.Context, organizationID string, filters ListFilters) (*ListResult, error) {
	page, pageSize := normalizePage(filters.Page, filters.PageSize)

	where := []string{"organization_id = $1"}
	args := []any{organizationID}
	argPos := 2

	if filters.Action != "" {
		where = append(where, fmt.Sprintf("action = $%d", argPos))
		args = append(args, filters.Action)
		argPos++
	}
	if filters.EntityType != "" {
		where = append(where, fmt.Sprintf("entity_type = $%d", argPos))
		args = append(args, filters.EntityType)
		argPos++
	}
	if filters.EntityID != "" {
		where = append(where, fmt.Sprintf("entity_id = $%d", argPos))
		args = append(args, filters.EntityID)
		argPos++
	}
	if filters.FromDate != "" {
		where = append(where, fmt.Sprintf("created_at >= $%d::timestamptz", argPos))
		args = append(args, filters.FromDate)
		argPos++
	}
	if filters.ToDate != "" {
		where = append(where, fmt.Sprintf("created_at <= $%d::timestamptz", argPos))
		args = append(args, filters.ToDate)
		argPos++
	}
	if filters.ActorUserID != "" {
		where = append(where, fmt.Sprintf("actor_user_id = $%d", argPos))
		args = append(args, filters.ActorUserID)
		argPos++
	}

	limitPos := argPos
	offsetPos := argPos + 1
	args = append(args, pageSize, (page-1)*pageSize)

	query := fmt.Sprintf(`
		SELECT
			id::text,
			organization_id::text,
			actor_user_id::text,
			action,
			entity_type,
			entity_id::text,
			old_value,
			new_value,
			ip_address,
			user_agent,
			created_at
		FROM audit_logs
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, strings.Join(where, " AND "), limitPos, offsetPos)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]AuditLog, 0)
	for rows.Next() {
		var item AuditLog
		var actorUserID, entityID, ipAddress, userAgent sql.NullString
		var oldValue, newValue json.RawMessage
		if err := rows.Scan(&item.ID, &item.OrganizationID, &actorUserID, &item.Action, &item.EntityType, &entityID, &oldValue, &newValue, &ipAddress, &userAgent, &item.CreatedAt); err != nil {
			return nil, err
		}
		if actorUserID.Valid {
			value := actorUserID.String
			item.ActorUserID = &value
		}
		if entityID.Valid {
			value := entityID.String
			item.EntityID = &value
		}
		if len(oldValue) > 0 {
			item.OldValue = append(json.RawMessage(nil), oldValue...)
		}
		if len(newValue) > 0 {
			item.NewValue = append(json.RawMessage(nil), newValue...)
		}
		if ipAddress.Valid {
			value := ipAddress.String
			item.IPAddr = &value
		}
		if userAgent.Valid {
			value := userAgent.String
			item.UserAgent = &value
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var total int
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM audit_logs WHERE %s`, strings.Join(where, " AND "))
	if err := s.db.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total); err != nil {
		return nil, err
	}

	return &ListResult{Items: items, Pagination: Pagination{Page: page, PageSize: pageSize, Total: total}}, nil
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

var ErrNotFound = errors.New("not found")

func (s *Service) record(ctx context.Context, runner interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}, organizationID string, actorUserID *string, action, entityType string, entityID *string, oldValue, newValue any, ipAddress, userAgent *string) error {
	var oldJSON, newJSON any
	if oldValue != nil {
		payload, err := json.Marshal(oldValue)
		if err != nil {
			return err
		}
		oldJSON = payload
	}
	if newValue != nil {
		payload, err := json.Marshal(newValue)
		if err != nil {
			return err
		}
		newJSON = payload
	}

	_, err := runner.ExecContext(ctx, `
		INSERT INTO audit_logs (
			organization_id, actor_user_id, action, entity_type, entity_id, old_value, new_value, ip_address, user_agent
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, organizationID, actorUserID, action, entityType, entityID, oldJSON, newJSON, ipAddress, userAgent)
	return err
}
