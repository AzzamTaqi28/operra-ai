package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"

	"operra/api/internal/platform/security"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

type User struct {
	ID             string   `json:"id"`
	OrganizationID string   `json:"organization_id"`
	DepartmentID   *string  `json:"department_id,omitempty"`
	Name           string   `json:"name"`
	Email          string   `json:"email"`
	Status         string   `json:"status"`
	Roles          []string `json:"roles"`
}

type ListFilters struct {
	Page         int
	PageSize     int
	Status       string
	Role         string
	DepartmentID string
	Search       string
}

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

type ListResult struct {
	Items      []User     `json:"items"`
	Pagination Pagination `json:"pagination"`
}

type CreateRequest struct {
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	Password     string   `json:"password"`
	DepartmentID string   `json:"department_id"`
	RoleKeys     []string `json:"role_keys"`
}

type UpdateRequest struct {
	Name         *string `json:"name"`
	Email        *string `json:"email"`
	Password     *string `json:"password"`
	DepartmentID *string `json:"department_id"`
	Status       *string `json:"status"`
}

func (s *Service) List(ctx context.Context, organizationID string, filters ListFilters) (*ListResult, error) {
	page, pageSize := normalizePage(filters.Page, filters.PageSize)

	where := []string{"organization_id = $1"}
	args := []any{organizationID}
	argPos := 2

	if filters.Status != "" {
		where = append(where, fmt.Sprintf("status = $%d", argPos))
		args = append(args, filters.Status)
		argPos++
	}

	if filters.DepartmentID != "" {
		where = append(where, fmt.Sprintf("department_id = $%d", argPos))
		args = append(args, filters.DepartmentID)
		argPos++
	}

	if filters.Search != "" {
		where = append(where, fmt.Sprintf("(name ILIKE $%d OR email ILIKE $%d)", argPos, argPos+1))
		term := "%" + filters.Search + "%"
		args = append(args, term, term)
		argPos += 2
	}

	if filters.Role != "" {
		where = append(where, fmt.Sprintf(`EXISTS (
			SELECT 1
			FROM user_roles ur
			INNER JOIN roles r ON r.id = ur.role_id
			WHERE ur.organization_id = users.organization_id
			  AND ur.user_id = users.id
			  AND r.key = $%d
		)`, argPos))
		args = append(args, filters.Role)
		argPos++
	}

	whereSQL := strings.Join(where, " AND ")
	limitPos := argPos
	offsetPos := argPos + 1
	args = append(args, pageSize, (page-1)*pageSize)

	query := fmt.Sprintf(`
		SELECT id::text, organization_id::text, department_id::text, name, email, status
		FROM users
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereSQL, limitPos, offsetPos)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]User, 0)
	userIDs := make([]string, 0)
	for rows.Next() {
		var user User
		var departmentID sql.NullString
		if err := rows.Scan(&user.ID, &user.OrganizationID, &departmentID, &user.Name, &user.Email, &user.Status); err != nil {
			return nil, err
		}
		if departmentID.Valid {
			value := departmentID.String
			user.DepartmentID = &value
		}
		items = append(items, user)
		userIDs = append(userIDs, user.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	rolesByUser, err := s.rolesByUserIDs(ctx, organizationID, userIDs)
	if err != nil {
		return nil, err
	}

	for idx := range items {
		items[idx].Roles = rolesByUser[items[idx].ID]
	}

	total, err := s.count(ctx, organizationID, whereSQL, args[:len(args)-2])
	if err != nil {
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

func (s *Service) Get(ctx context.Context, organizationID, userID string) (*User, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id::text, organization_id::text, department_id::text, name, email, status
		FROM users
		WHERE organization_id = $1 AND id = $2
	`, organizationID, userID)

	var user User
	var departmentID sql.NullString
	if err := row.Scan(&user.ID, &user.OrganizationID, &departmentID, &user.Name, &user.Email, &user.Status); err != nil {
		return nil, err
	}
	if departmentID.Valid {
		value := departmentID.String
		user.DepartmentID = &value
	}

	roles, err := s.rolesByUserIDs(ctx, organizationID, []string{userID})
	if err != nil {
		return nil, err
	}
	user.Roles = roles[userID]
	return &user, nil
}

func (s *Service) Create(ctx context.Context, organizationID string, req CreateRequest) (*User, error) {
	if err := validateCreate(req); err != nil {
		return nil, err
	}

	passwordHash, err := security.HashPassword(strings.TrimSpace(req.Password))
	if err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	departmentID := sql.NullString{}
	if strings.TrimSpace(req.DepartmentID) != "" {
		departmentID = sql.NullString{String: strings.TrimSpace(req.DepartmentID), Valid: true}
		if err := verifyDepartment(ctx, tx, organizationID, departmentID.String); err != nil {
			return nil, err
		}
	}

	var userID string
	err = tx.QueryRowContext(ctx, `
		INSERT INTO users (organization_id, department_id, name, email, password_hash)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id::text
	`, organizationID, nullableString(departmentID), strings.TrimSpace(req.Name), normalizeEmail(req.Email), passwordHash).Scan(&userID)
	if err != nil {
		return nil, translateError(err)
	}

	if err := s.assignRolesTx(ctx, tx, organizationID, userID, req.RoleKeys); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.Get(ctx, organizationID, userID)
}

func (s *Service) Update(ctx context.Context, organizationID, userID string, req UpdateRequest) (*User, error) {
	if err := validateUpdate(req); err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	current, err := s.Get(ctx, organizationID, userID)
	if err != nil {
		return nil, err
	}

	name := current.Name
	if req.Name != nil {
		name = strings.TrimSpace(*req.Name)
	}

	email := current.Email
	if req.Email != nil {
		email = normalizeEmail(*req.Email)
	}

	status := current.Status
	if req.Status != nil {
		status = strings.TrimSpace(*req.Status)
	}

	departmentID := sql.NullString{}
	if current.DepartmentID != nil {
		departmentID = sql.NullString{String: *current.DepartmentID, Valid: true}
	}
	if req.DepartmentID != nil {
		if strings.TrimSpace(*req.DepartmentID) == "" {
			departmentID = sql.NullString{}
		} else {
			departmentID = sql.NullString{String: strings.TrimSpace(*req.DepartmentID), Valid: true}
			if err := verifyDepartment(ctx, tx, organizationID, departmentID.String); err != nil {
				return nil, err
			}
		}
	}

	if req.Password != nil && strings.TrimSpace(*req.Password) != "" {
		passwordHash, err := security.HashPassword(strings.TrimSpace(*req.Password))
		if err != nil {
			return nil, err
		}
		if _, err := tx.ExecContext(ctx, `
			UPDATE users
			SET department_id = $3,
			    name = $4,
			    email = $5,
			    password_hash = $6,
			    status = $7,
			    updated_at = NOW()
			WHERE organization_id = $1 AND id = $2
		`, organizationID, userID, nullableString(departmentID), name, email, passwordHash, status); err != nil {
			return nil, translateError(err)
		}
	} else {
		if _, err := tx.ExecContext(ctx, `
			UPDATE users
			SET department_id = $3,
			    name = $4,
			    email = $5,
			    status = $6,
			    updated_at = NOW()
			WHERE organization_id = $1 AND id = $2
		`, organizationID, userID, nullableString(departmentID), name, email, status); err != nil {
			return nil, translateError(err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.Get(ctx, organizationID, userID)
}

func (s *Service) AssignRoles(ctx context.Context, organizationID, userID string, roleKeys []string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if err := s.assignRolesTx(ctx, tx, organizationID, userID, roleKeys); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Service) RemoveRole(ctx context.Context, organizationID, userID, roleKey string) error {
	result, err := s.db.ExecContext(ctx, `
		DELETE FROM user_roles
		WHERE organization_id = $1
		  AND user_id = $2
		  AND role_id IN (
			SELECT id FROM roles WHERE organization_id = $1 AND key = $3
		)
	`, organizationID, userID, roleKey)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *Service) assignRolesTx(ctx context.Context, tx *sql.Tx, organizationID, userID string, roleKeys []string) error {
	if len(roleKeys) == 0 {
		return nil
	}

	normalized := normalizeRoleKeys(roleKeys)
	rows, err := tx.QueryContext(ctx, `
		SELECT key
		FROM roles
		WHERE organization_id = $1 AND key = ANY($2)
	`, organizationID, pq.Array(normalized))
	if err != nil {
		return err
	}
	defer rows.Close()

	found := make(map[string]struct{}, len(normalized))
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return err
		}
		found[key] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if len(found) != len(normalized) {
		return errors.New("one or more role_keys are invalid")
	}

	for _, roleKey := range normalized {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO user_roles (organization_id, user_id, role_id)
			SELECT $1, $2, id
			FROM roles
			WHERE organization_id = $1 AND key = $3
			ON CONFLICT (organization_id, user_id, role_id) DO NOTHING
		`, organizationID, userID, roleKey); err != nil {
			return translateError(err)
		}
	}

	return nil
}

func (s *Service) rolesByUserIDs(ctx context.Context, organizationID string, userIDs []string) (map[string][]string, error) {
	result := make(map[string][]string, len(userIDs))
	if len(userIDs) == 0 {
		return result, nil
	}

	rows, err := s.db.QueryContext(ctx, `
		SELECT ur.user_id::text, r.key
		FROM user_roles ur
		INNER JOIN roles r ON r.id = ur.role_id
		WHERE ur.organization_id = $1
		  AND ur.user_id = ANY($2)
		ORDER BY r.key
	`, organizationID, pq.Array(userIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID, key string
		if err := rows.Scan(&userID, &key); err != nil {
			return nil, err
		}
		result[userID] = append(result[userID], key)
	}

	return result, rows.Err()
}

func (s *Service) count(ctx context.Context, organizationID, whereSQL string, args []any) (int, error) {
	query := fmt.Sprintf(`SELECT COUNT(*) FROM users WHERE %s`, whereSQL)
	var total int
	if err := s.db.QueryRowContext(ctx, query, args...).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func validateCreate(req CreateRequest) error {
	switch {
	case strings.TrimSpace(req.Name) == "":
		return errors.New("name is required")
	case strings.TrimSpace(req.Email) == "":
		return errors.New("email is required")
	case strings.TrimSpace(req.Password) == "":
		return errors.New("password is required")
	}
	return nil
}

func validateUpdate(req UpdateRequest) error {
	if req.Name == nil && req.Email == nil && req.Password == nil && req.DepartmentID == nil && req.Status == nil {
		return errors.New("at least one field is required")
	}
	return nil
}

func translateError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		return ErrConflict
	}
	return err
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

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func normalizeRoleKeys(keys []string) []string {
	seen := make(map[string]struct{}, len(keys))
	normalized := make([]string, 0, len(keys))
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		normalized = append(normalized, key)
	}
	return normalized
}

func nullableString(value sql.NullString) any {
	if !value.Valid {
		return nil
	}
	return value.String
}

func verifyDepartment(ctx context.Context, db queryer, organizationID, departmentID string) error {
	var exists string
	if err := db.QueryRowContext(ctx, `
		SELECT id::text
		FROM departments
		WHERE organization_id = $1 AND id = $2
	`, organizationID, departmentID).Scan(&exists); err != nil {
		return err
	}
	return nil
}

type queryer interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

var ErrConflict = errors.New("conflict")
