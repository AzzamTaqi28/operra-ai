package departments

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/lib/pq"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

type Department struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
	Code           string `json:"code,omitempty"`
}

type CreateRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type UpdateRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func (s *Service) List(ctx context.Context, organizationID string) ([]Department, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id::text, organization_id::text, name, COALESCE(code, '')
		FROM departments
		WHERE organization_id = $1
		ORDER BY name ASC
	`, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Department
	for rows.Next() {
		var dept Department
		if err := rows.Scan(&dept.ID, &dept.OrganizationID, &dept.Name, &dept.Code); err != nil {
			return nil, err
		}
		items = append(items, dept)
	}

	return items, rows.Err()
}

func (s *Service) Create(ctx context.Context, organizationID string, req CreateRequest) (*Department, error) {
	if err := validateDepartment(req); err != nil {
		return nil, err
	}

	var dept Department
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO departments (organization_id, name, code)
		VALUES ($1, $2, NULLIF($3, ''))
		RETURNING id::text, organization_id::text, name, COALESCE(code, '')
	`, organizationID, strings.TrimSpace(req.Name), strings.TrimSpace(req.Code)).Scan(&dept.ID, &dept.OrganizationID, &dept.Name, &dept.Code)
	if err != nil {
		return nil, translateError(err)
	}

	return &dept, nil
}

func (s *Service) Update(ctx context.Context, organizationID, departmentID string, req UpdateRequest) (*Department, error) {
	if err := validateDepartment(req); err != nil {
		return nil, err
	}

	var dept Department
	err := s.db.QueryRowContext(ctx, `
		UPDATE departments
		SET name = $3,
		    code = NULLIF($4, ''),
		    updated_at = NOW()
		WHERE organization_id = $1 AND id = $2
		RETURNING id::text, organization_id::text, name, COALESCE(code, '')
	`, organizationID, departmentID, strings.TrimSpace(req.Name), strings.TrimSpace(req.Code)).Scan(&dept.ID, &dept.OrganizationID, &dept.Name, &dept.Code)
	if err != nil {
		return nil, translateError(err)
	}

	return &dept, nil
}

func (s *Service) Delete(ctx context.Context, organizationID, departmentID string) error {
	result, err := s.db.ExecContext(ctx, `
		DELETE FROM departments
		WHERE organization_id = $1 AND id = $2
	`, organizationID, departmentID)
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

func validateDepartment(req interface{ GetName() string }) error {
	if strings.TrimSpace(req.GetName()) == "" {
		return errors.New("name is required")
	}
	return nil
}

func (r CreateRequest) GetName() string { return r.Name }
func (r UpdateRequest) GetName() string { return r.Name }

func translateError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		return ErrConflict
	}
	return err
}

var ErrConflict = errors.New("conflict")
