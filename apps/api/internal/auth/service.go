package auth

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/lib/pq"

	"operra/api/internal/platform/middleware"
	"operra/api/internal/platform/security"
)

var builtinRoles = []struct {
	Key         string
	Name        string
	Description string
}{
	{Key: "owner", Name: "Owner", Description: "Organization owner"},
	{Key: "admin", Name: "Admin", Description: "Internal admin"},
	{Key: "requester", Name: "Requester", Description: "Request creator"},
	{Key: "manager", Name: "Manager", Description: "Department approver"},
	{Key: "finance", Name: "Finance", Description: "Finance approver"},
	{Key: "procurement", Name: "Procurement", Description: "Procurement processor"},
	{Key: "director", Name: "Director", Description: "High-value approver"},
	{Key: "auditor", Name: "Auditor", Description: "Read-only auditor"},
}

type Service struct {
	db     *sql.DB
	secret string
}

func NewService(db *sql.DB, secret string) *Service {
	return &Service{db: db, secret: secret}
}

type RegisterOrganizationRequest struct {
	OrganizationName string
	OrganizationSlug string
	OwnerName        string
	OwnerEmail       string
	Password         string
}

type LoginRequest struct {
	Email    string
	Password string
}

type OrganizationResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type UserResponse struct {
	ID             string   `json:"id"`
	OrganizationID string   `json:"organization_id,omitempty"`
	Name           string   `json:"name"`
	Email          string   `json:"email"`
	Roles          []string `json:"roles,omitempty"`
}

type RegisterResponse struct {
	Organization OrganizationResponse `json:"organization"`
	User         UserResponse         `json:"user"`
	Token        string               `json:"token"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type CurrentUserResponse = middleware.CurrentUser

func (s *Service) RegisterOrganization(ctx context.Context, req RegisterOrganizationRequest) (*RegisterResponse, error) {
	if err := validateRegistration(req); err != nil {
		return nil, err
	}

	passwordHash, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var organizationID, userID string
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO organizations (name, slug)
		VALUES ($1, $2)
		RETURNING id::text
	`, req.OrganizationName, req.OrganizationSlug).Scan(&organizationID); err != nil {
		return nil, translateDBError(err)
	}

	if err := tx.QueryRowContext(ctx, `
		INSERT INTO users (organization_id, name, email, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id::text
	`, organizationID, req.OwnerName, strings.ToLower(strings.TrimSpace(req.OwnerEmail)), passwordHash).Scan(&userID); err != nil {
		return nil, translateDBError(err)
	}

	for _, role := range builtinRoles {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO roles (organization_id, key, name, description, is_system)
			VALUES ($1, $2, $3, $4, true)
			ON CONFLICT (organization_id, key) DO NOTHING
		`, organizationID, role.Key, role.Name, role.Description); err != nil {
			return nil, translateDBError(err)
		}
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO user_roles (organization_id, user_id, role_id)
		SELECT $1, $2, id
		FROM roles
		WHERE organization_id = $1 AND key = 'owner'
	`, organizationID, userID); err != nil {
		return nil, translateDBError(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	token, err := security.GenerateToken(s.secret, security.TokenClaims{
		UserID:         userID,
		OrganizationID: organizationID,
		Roles:          []string{"owner"},
	}, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &RegisterResponse{
		Organization: OrganizationResponse{
			ID:   organizationID,
			Name: req.OrganizationName,
			Slug: req.OrganizationSlug,
		},
		User: UserResponse{
			ID:             userID,
			OrganizationID: organizationID,
			Name:           req.OwnerName,
			Email:          strings.ToLower(strings.TrimSpace(req.OwnerEmail)),
			Roles:          []string{"owner"},
		},
		Token: token,
	}, nil
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}

	row := s.db.QueryRowContext(ctx, `
		SELECT
			u.id::text,
			u.organization_id::text,
			u.name,
			u.email,
			u.password_hash,
			u.status
		FROM users u
		WHERE lower(u.email) = lower($1)
		ORDER BY u.created_at ASC
		LIMIT 1
	`, email)

	var userID, organizationID, name, dbEmail, passwordHash, status string
	if err := row.Scan(&userID, &organizationID, &name, &dbEmail, &passwordHash, &status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if status != "active" {
		return nil, errors.New("invalid credentials")
	}

	if err := security.ComparePassword(passwordHash, req.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	roles, err := s.rolesForUser(ctx, organizationID, userID)
	if err != nil {
		return nil, err
	}

	token, err := security.GenerateToken(s.secret, security.TokenClaims{
		UserID:         userID,
		OrganizationID: organizationID,
		Roles:          roles,
	}, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User: UserResponse{
			ID:             userID,
			OrganizationID: organizationID,
			Name:           name,
			Email:          dbEmail,
			Roles:          roles,
		},
	}, nil
}

func (s *Service) CurrentUser(userID, organizationID string) (*middleware.CurrentUser, error) {
	row := s.db.QueryRowContext(context.Background(), `
		SELECT id::text, organization_id::text, department_id::text, name, email, status
		FROM users
		WHERE id = $1 AND organization_id = $2
	`, userID, organizationID)

	var user middleware.CurrentUser
	var departmentID sql.NullString
	if err := row.Scan(&user.ID, &user.OrganizationID, &departmentID, &user.Name, &user.Email, &user.Status); err != nil {
		return nil, err
	}
	if departmentID.Valid {
		value := departmentID.String
		user.DepartmentID = &value
	}

	roles, err := s.rolesForUser(context.Background(), organizationID, userID)
	if err != nil {
		return nil, err
	}
	user.Roles = roles
	return &user, nil
}

func (s *Service) rolesForUser(ctx context.Context, organizationID, userID string) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT r.key
		FROM roles r
		INNER JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.organization_id = $1 AND ur.user_id = $2
		ORDER BY r.key
	`, organizationID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, rows.Err()
}

func validateRegistration(req RegisterOrganizationRequest) error {
	switch {
	case strings.TrimSpace(req.OrganizationName) == "":
		return errors.New("organization_name is required")
	case strings.TrimSpace(req.OrganizationSlug) == "":
		return errors.New("organization_slug is required")
	case strings.TrimSpace(req.OwnerName) == "":
		return errors.New("owner_name is required")
	case strings.TrimSpace(req.OwnerEmail) == "":
		return errors.New("owner_email is required")
	case strings.TrimSpace(req.Password) == "":
		return errors.New("password is required")
	}

	return nil
}

func translateDBError(err error) error {
	if err == nil {
		return nil
	}
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		return ErrConflict
	}
	return err
}

var ErrConflict = errors.New("conflict")
