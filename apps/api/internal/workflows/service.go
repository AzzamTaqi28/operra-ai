package workflows

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"

	configworkflow "operra/api/internal/workflow"
)

type MermaidGenerator func(configworkflow.Config) string
type ConfigValidator func(configworkflow.Config) configworkflow.ValidationResult

type Service struct {
	db       *sql.DB
	generate MermaidGenerator
	validate ConfigValidator
}

func NewService(db *sql.DB, generate MermaidGenerator, validate ConfigValidator) *Service {
	return &Service{db: db, generate: generate, validate: validate}
}

type Workflow struct {
	ID              string    `json:"id"`
	OrganizationID  string    `json:"organization_id"`
	Name            string    `json:"name"`
	Type            string    `json:"type"`
	Status          string    `json:"status"`
	ActiveVersionID *string   `json:"active_version_id,omitempty"`
	CreatedBy       string    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	ActiveVersion   *Version  `json:"active_version,omitempty"`
}

type Version struct {
	ID             string                `json:"id"`
	OrganizationID string                `json:"organization_id"`
	WorkflowID     string                `json:"workflow_id"`
	VersionNumber  int                   `json:"version_number"`
	ConfigJSON     configworkflow.Config `json:"config_json"`
	MermaidDiagram string                `json:"mermaid_diagram,omitempty"`
	Explanation    string                `json:"explanation,omitempty"`
	CreatedBy      string                `json:"created_by"`
	CreatedAt      time.Time             `json:"created_at"`
}

type VersionSummary struct {
	ID             string                `json:"id"`
	WorkflowID     string                `json:"workflow_id"`
	VersionNumber  int                   `json:"version_number"`
	ConfigJSON     configworkflow.Config `json:"config_json"`
	MermaidDiagram string                `json:"mermaid_diagram,omitempty"`
	Explanation    string                `json:"explanation,omitempty"`
	CreatedBy      string                `json:"created_by"`
	CreatedAt      time.Time             `json:"created_at"`
}

type CreateRequest struct {
	Name       string          `json:"name"`
	Type       string          `json:"type"`
	ConfigJSON json.RawMessage `json:"config_json"`
}

type CreateVersionRequest struct {
	ConfigJSON json.RawMessage `json:"config_json"`
}

type ValidationResult = configworkflow.ValidationResult
type ValidationError = configworkflow.ValidationError

type CreateResult struct {
	Workflow Workflow `json:"workflow"`
	Version  Version  `json:"version"`
}

type VersionCreateResult struct {
	Version    Version          `json:"version"`
	Validation ValidationResult `json:"validation"`
}

func (s *Service) List(ctx context.Context, organizationID, filterType, status string) ([]Workflow, error) {
	args := []any{organizationID}
	where := []string{"w.organization_id = $1"}
	argPos := 2

	if filterType != "" {
		where = append(where, fmt.Sprintf("w.type = $%d", argPos))
		args = append(args, filterType)
		argPos++
	}
	if status != "" {
		where = append(where, fmt.Sprintf("w.status = $%d", argPos))
		args = append(args, status)
		argPos++
	}

	query := fmt.Sprintf(`
		SELECT
			w.id::text,
			w.organization_id::text,
			w.name,
			w.type,
			w.status,
			w.active_version_id::text,
			w.created_by::text,
			w.created_at,
			w.updated_at,
			av.id::text,
			av.organization_id::text,
			av.workflow_id::text,
			av.version_number,
			av.config_json,
			COALESCE(av.mermaid_diagram, ''),
			COALESCE(av.explanation, ''),
			av.created_by::text,
			av.created_at
		FROM workflows w
		LEFT JOIN workflow_versions av ON av.id = w.active_version_id
		WHERE %s
		ORDER BY w.created_at DESC
	`, strings.Join(where, " AND "))

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workflows []Workflow
	for rows.Next() {
		var workflow Workflow
		var activeVersionID sql.NullString
		var activeVersion Version
		var activeVersionIDScan sql.NullString
		var avOrganizationID sql.NullString
		var avWorkflowID sql.NullString
		var avVersionNumber sql.NullInt64
		var configJSONBytes []byte
		var mermaid, explanation sql.NullString
		var avCreatedBy sql.NullString
		var avCreatedAt sql.NullTime
		if err := rows.Scan(
			&workflow.ID,
			&workflow.OrganizationID,
			&workflow.Name,
			&workflow.Type,
			&workflow.Status,
			&activeVersionID,
			&workflow.CreatedBy,
			&workflow.CreatedAt,
			&workflow.UpdatedAt,
			&activeVersionIDScan,
			&avOrganizationID,
			&avWorkflowID,
			&avVersionNumber,
			&configJSONBytes,
			&mermaid,
			&explanation,
			&avCreatedBy,
			&avCreatedAt,
		); err != nil {
			return nil, err
		}
		if activeVersionID.Valid {
			value := activeVersionID.String
			workflow.ActiveVersionID = &value
		}
		if activeVersionIDScan.Valid {
			activeVersion.ID = activeVersionIDScan.String
			if avOrganizationID.Valid {
				activeVersion.OrganizationID = avOrganizationID.String
			}
			if avWorkflowID.Valid {
				activeVersion.WorkflowID = avWorkflowID.String
			}
			if avVersionNumber.Valid {
				activeVersion.VersionNumber = int(avVersionNumber.Int64)
			}
			if mermaid.Valid {
				activeVersion.MermaidDiagram = mermaid.String
			}
			if explanation.Valid {
				activeVersion.Explanation = explanation.String
			}
			if avCreatedBy.Valid {
				activeVersion.CreatedBy = avCreatedBy.String
			}
			if avCreatedAt.Valid {
				activeVersion.CreatedAt = avCreatedAt.Time
			}
			if len(configJSONBytes) > 0 {
				var cfg configworkflow.Config
				if err := json.Unmarshal(configJSONBytes, &cfg); err != nil {
					return nil, err
				}
				activeVersion.ConfigJSON = cfg
			}
			workflow.ActiveVersion = &activeVersion
		}
		workflows = append(workflows, workflow)
	}
	return workflows, rows.Err()
}

func (s *Service) Get(ctx context.Context, organizationID, workflowID string) (*Workflow, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT
			w.id::text,
			w.organization_id::text,
			w.name,
			w.type,
			w.status,
			w.active_version_id::text,
			w.created_by::text,
			w.created_at,
			w.updated_at,
			av.id::text,
			av.organization_id::text,
			av.workflow_id::text,
			av.version_number,
			av.config_json,
			av.mermaid_diagram,
			av.explanation,
			av.created_by::text,
			av.created_at
		FROM workflows w
		LEFT JOIN workflow_versions av ON av.id = w.active_version_id
		WHERE w.organization_id = $1 AND w.id = $2
	`, organizationID, workflowID)

	return scanWorkflowRow(row)
}

func (s *Service) ListVersions(ctx context.Context, organizationID, workflowID string) ([]VersionSummary, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT
			id::text,
			workflow_id::text,
			version_number,
			config_json,
			COALESCE(mermaid_diagram, ''),
			COALESCE(explanation, ''),
			created_by::text,
			created_at
		FROM workflow_versions
		WHERE organization_id = $1 AND workflow_id = $2
		ORDER BY version_number DESC
	`, organizationID, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []VersionSummary
	for rows.Next() {
		var version VersionSummary
		var cfgBytes []byte
		if err := rows.Scan(
			&version.ID,
			&version.WorkflowID,
			&version.VersionNumber,
			&cfgBytes,
			&version.MermaidDiagram,
			&version.Explanation,
			&version.CreatedBy,
			&version.CreatedAt,
		); err != nil {
			return nil, err
		}
		if len(cfgBytes) > 0 {
			if err := json.Unmarshal(cfgBytes, &version.ConfigJSON); err != nil {
				return nil, err
			}
		}
		versions = append(versions, version)
	}
	return versions, rows.Err()
}

func (s *Service) CreateWorkflow(ctx context.Context, organizationID, createdBy string, req CreateRequest) (*CreateResult, error) {
	cfg, validation, err := s.buildConfig(req.Name, req.Type, req.ConfigJSON, 1)
	if err != nil {
		return nil, err
	}
	if !validation.Valid {
		return nil, ErrInvalidWorkflow{Validation: validation}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	var workflowID string
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO workflows (organization_id, name, type, status, created_by)
		VALUES ($1, $2, $3, 'draft', $4)
		RETURNING id::text
	`, organizationID, cfg.Name, cfg.Type, createdBy).Scan(&workflowID); err != nil {
		return nil, translateError(err)
	}

	version, err := s.insertVersion(ctx, tx, organizationID, workflowID, createdBy, cfg, 1, "")
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &CreateResult{
		Workflow: Workflow{
			ID:             workflowID,
			OrganizationID: organizationID,
			Name:           cfg.Name,
			Type:           cfg.Type,
			Status:         "draft",
			CreatedBy:      createdBy,
			CreatedAt:      version.CreatedAt,
			UpdatedAt:      version.CreatedAt,
		},
		Version: version,
	}, nil
}

func (s *Service) CreateVersion(ctx context.Context, organizationID, createdBy, workflowID string, req CreateVersionRequest) (*VersionCreateResult, error) {
	workflow, err := s.Get(ctx, organizationID, workflowID)
	if err != nil {
		return nil, err
	}

	nextVersion, err := s.nextVersionNumber(ctx, organizationID, workflowID)
	if err != nil {
		return nil, err
	}

	cfg, validation, err := s.buildConfig(workflow.Name, workflow.Type, req.ConfigJSON, nextVersion)
	if err != nil {
		return nil, err
	}
	if !validation.Valid {
		return nil, ErrInvalidWorkflow{Validation: validation}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	version, err := s.insertVersion(ctx, tx, organizationID, workflowID, createdBy, cfg, nextVersion, "")
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &VersionCreateResult{
		Version: version,
		Validation: ValidationResult{
			Valid:  true,
			Errors: nil,
		},
	}, nil
}

func (s *Service) ActivateVersion(ctx context.Context, organizationID, workflowID, versionID string) error {
	result, err := s.db.ExecContext(ctx, `
		UPDATE workflows
		SET active_version_id = $3,
		    status = 'active',
		    updated_at = NOW()
		WHERE organization_id = $1
		  AND id = $2
		  AND EXISTS (
			SELECT 1
			FROM workflow_versions
			WHERE organization_id = $1
			  AND workflow_id = $2
			  AND id = $3
		  )
	`, organizationID, workflowID, versionID)
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

func (s *Service) buildConfig(name, typ string, raw json.RawMessage, version int) (configworkflow.Config, configworkflow.ValidationResult, error) {
	var cfg configworkflow.Config
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &cfg); err != nil {
			return configworkflow.Config{}, configworkflow.ValidationResult{}, err
		}
	}
	if name != "" {
		cfg.Name = name
	}
	if typ != "" {
		cfg.Type = typ
	}
	if cfg.Version == 0 {
		cfg.Version = version
	}
	if len(cfg.Steps) == 0 && strings.TrimSpace(cfg.Name) == "" {
		// Let validator surface the detailed errors.
	}

	validation := s.validate(cfg)
	if cfg.Version != version {
		cfg.Version = version
	}
	return cfg, validation, nil
}

func (s *Service) insertVersion(ctx context.Context, tx *sql.Tx, organizationID, workflowID, createdBy string, cfg configworkflow.Config, versionNumber int, explanation string) (Version, error) {
	mermaid := s.generate(cfg)
	payload, err := json.Marshal(cfg)
	if err != nil {
		return Version{}, err
	}

	var version Version
	err = tx.QueryRowContext(ctx, `
		INSERT INTO workflow_versions (
			organization_id, workflow_id, version_number, config_json, mermaid_diagram, explanation, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING
			id::text,
			organization_id::text,
			workflow_id::text,
			version_number,
			created_by::text,
			created_at
	`, organizationID, workflowID, versionNumber, payload, mermaid, explanation, createdBy).Scan(
		&version.ID,
		&version.OrganizationID,
		&version.WorkflowID,
		&version.VersionNumber,
		&version.CreatedBy,
		&version.CreatedAt,
	)
	if err != nil {
		return Version{}, translateError(err)
	}
	version.ConfigJSON = cfg
	version.MermaidDiagram = mermaid
	version.Explanation = explanation
	return version, nil
}

func scanWorkflowRow(scanner interface {
	Scan(dest ...any) error
}) (*Workflow, error) {
	var workflow Workflow
	var activeVersionID sql.NullString
	var activeVersion Version
	var activeVersionIDScan sql.NullString
	var avOrganizationID sql.NullString
	var avWorkflowID sql.NullString
	var avVersionNumber sql.NullInt64
	var configJSONBytes []byte
	var mermaid, explanation sql.NullString
	var avCreatedBy sql.NullString
	var avCreatedAt sql.NullTime
	if err := scanner.Scan(
		&workflow.ID,
		&workflow.OrganizationID,
		&workflow.Name,
		&workflow.Type,
		&workflow.Status,
		&activeVersionID,
		&workflow.CreatedBy,
		&workflow.CreatedAt,
		&workflow.UpdatedAt,
		&activeVersionIDScan,
		&avOrganizationID,
		&avWorkflowID,
		&avVersionNumber,
		&configJSONBytes,
		&mermaid,
		&explanation,
		&avCreatedBy,
		&avCreatedAt,
	); err != nil {
		return nil, err
	}
	if activeVersionID.Valid {
		value := activeVersionID.String
		workflow.ActiveVersionID = &value
	}
	if activeVersionIDScan.Valid {
		activeVersion.ID = activeVersionIDScan.String
		if avOrganizationID.Valid {
			activeVersion.OrganizationID = avOrganizationID.String
		}
		if avWorkflowID.Valid {
			activeVersion.WorkflowID = avWorkflowID.String
		}
		if avVersionNumber.Valid {
			activeVersion.VersionNumber = int(avVersionNumber.Int64)
		}
		if len(configJSONBytes) > 0 {
			var cfg configworkflow.Config
			if err := json.Unmarshal(configJSONBytes, &cfg); err != nil {
				return nil, err
			}
			activeVersion.ConfigJSON = cfg
		}
		if mermaid.Valid {
			activeVersion.MermaidDiagram = mermaid.String
		}
		if explanation.Valid {
			activeVersion.Explanation = explanation.String
		}
		if avCreatedBy.Valid {
			activeVersion.CreatedBy = avCreatedBy.String
		}
		if avCreatedAt.Valid {
			activeVersion.CreatedAt = avCreatedAt.Time
		}
		workflow.ActiveVersion = &activeVersion
	}
	return &workflow, nil
}

func (s *Service) nextVersionNumber(ctx context.Context, organizationID, workflowID string) (int, error) {
	var next int
	if err := s.db.QueryRowContext(ctx, `
		SELECT COALESCE(MAX(version_number), 0) + 1
		FROM workflow_versions
		WHERE organization_id = $1 AND workflow_id = $2
	`, organizationID, workflowID).Scan(&next); err != nil {
		return 0, err
	}
	return next, nil
}

func translateError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		return ErrConflict
	}
	return err
}

type ErrInvalidWorkflow struct {
	Validation ValidationResult
}

func (e ErrInvalidWorkflow) Error() string {
	return "workflow invalid"
}

var ErrConflict = errors.New("conflict")
