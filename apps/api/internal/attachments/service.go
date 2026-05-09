package attachments

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"operra/api/internal/audit"
	"operra/api/internal/storage"
)

type Service struct {
	db    *sql.DB
	store storage.Store
	audit *audit.Service
}

var ErrForbidden = errors.New("forbidden")

func NewService(db *sql.DB, store storage.Store, auditSvc ...*audit.Service) *Service {
	var svc *audit.Service
	if len(auditSvc) > 0 {
		svc = auditSvc[0]
	}
	return &Service{db: db, store: store, audit: svc}
}

type Attachment struct {
	ID                string    `json:"id"`
	OrganizationID    string    `json:"organization_id"`
	PurchaseRequestID string    `json:"purchase_request_id"`
	UploadedBy        string    `json:"uploaded_by"`
	FileName          string    `json:"file_name"`
	FileSize          int64     `json:"file_size"`
	MIMEType          string    `json:"mime_type"`
	StorageDriver     string    `json:"storage_driver"`
	StorageBucket     string    `json:"storage_bucket"`
	StorageKey        string    `json:"storage_key"`
	Checksum          *string   `json:"checksum,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

type DownloadResult struct {
	Attachment Attachment
	Data       []byte
}

func (s *Service) Upload(ctx context.Context, organizationID string, actorID string, roles []string, requestID string, fileName, mimeType string, data []byte) (*Attachment, error) {
	current, err := s.requestAccess(ctx, organizationID, requestID)
	if err != nil {
		return nil, err
	}
	if current.RequesterID != actorID && !hasAnyRole(roles, "owner", "admin", "finance", "auditor", "manager", "director", "procurement") {
		return nil, ErrForbidden
	}
	if len(data) == 0 {
		return nil, errors.New("file is required")
	}

	attachmentID := fmt.Sprintf("%d", time.Now().UTC().UnixNano())
	safeName := sanitizeName(fileName)
	key := filepath.Join("attachments", organizationID, requestID, attachmentID+"_"+safeName)
	checksum := sha256Hex(data)

	if err := s.store.Put(ctx, key, mimeType, data); err != nil {
		return nil, err
	}

	row := s.db.QueryRowContext(ctx, `
		INSERT INTO attachments (
			organization_id, purchase_request_id, uploaded_by, file_name, file_size, mime_type, storage_driver, storage_bucket, storage_key, checksum
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id::text, organization_id::text, purchase_request_id::text, uploaded_by::text, file_name, file_size, mime_type, storage_driver, storage_bucket, storage_key, checksum::text, created_at
	`, organizationID, requestID, actorID, fileName, int64(len(data)), mimeType, storageDriverName(s.store), storageBucketName(s.store), key, checksum)

	var attachment Attachment
	var checksumValue sql.NullString
	if err := row.Scan(&attachment.ID, &attachment.OrganizationID, &attachment.PurchaseRequestID, &attachment.UploadedBy, &attachment.FileName, &attachment.FileSize, &attachment.MIMEType, &attachment.StorageDriver, &attachment.StorageBucket, &attachment.StorageKey, &checksumValue, &attachment.CreatedAt); err != nil {
		return nil, err
	}
	if checksumValue.Valid {
		value := checksumValue.String
		attachment.Checksum = &value
	}

	if err := s.recordAudit(ctx, organizationID, actorID, "attachment.uploaded", "attachment", &attachment.ID, nil, attachment); err != nil {
		return nil, err
	}

	return &attachment, nil
}

func (s *Service) Download(ctx context.Context, organizationID, actorID string, roles []string, requestID, attachmentID string) (*DownloadResult, error) {
	attachment, err := s.getAttachment(ctx, organizationID, requestID, attachmentID)
	if err != nil {
		return nil, err
	}
	if attachment.UploadedBy != actorID && !hasAnyRole(roles, "owner", "admin", "finance", "auditor", "manager", "director", "procurement") {
		return nil, ErrForbidden
	}

	data, err := s.store.Get(ctx, attachment.StorageKey)
	if err != nil {
		return nil, err
	}

	return &DownloadResult{Attachment: *attachment, Data: data}, nil
}

func (s *Service) ListByRequest(ctx context.Context, organizationID, requestID string) ([]Attachment, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id::text, organization_id::text, purchase_request_id::text, uploaded_by::text, file_name, file_size, mime_type, storage_driver, storage_bucket, storage_key, checksum::text, created_at
		FROM attachments
		WHERE organization_id = $1 AND purchase_request_id = $2
		ORDER BY created_at DESC
	`, organizationID, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Attachment
	for rows.Next() {
		var item Attachment
		var checksum sql.NullString
		if err := rows.Scan(&item.ID, &item.OrganizationID, &item.PurchaseRequestID, &item.UploadedBy, &item.FileName, &item.FileSize, &item.MIMEType, &item.StorageDriver, &item.StorageBucket, &item.StorageKey, &checksum, &item.CreatedAt); err != nil {
			return nil, err
		}
		if checksum.Valid {
			value := checksum.String
			item.Checksum = &value
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Service) requestAccess(ctx context.Context, organizationID, requestID string) (*requestRow, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id::text, organization_id::text, requester_id::text
		FROM purchase_requests
		WHERE organization_id = $1 AND id = $2
	`, organizationID, requestID)

	var item requestRow
	if err := row.Scan(&item.ID, &item.OrganizationID, &item.RequesterID); err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *Service) getAttachment(ctx context.Context, organizationID, requestID, attachmentID string) (*Attachment, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id::text, organization_id::text, purchase_request_id::text, uploaded_by::text, file_name, file_size, mime_type, storage_driver, storage_bucket, storage_key, checksum::text, created_at
		FROM attachments
		WHERE organization_id = $1 AND purchase_request_id = $2 AND id = $3
	`, organizationID, requestID, attachmentID)

	var attachment Attachment
	var checksum sql.NullString
	if err := row.Scan(&attachment.ID, &attachment.OrganizationID, &attachment.PurchaseRequestID, &attachment.UploadedBy, &attachment.FileName, &attachment.FileSize, &attachment.MIMEType, &attachment.StorageDriver, &attachment.StorageBucket, &attachment.StorageKey, &checksum, &attachment.CreatedAt); err != nil {
		return nil, err
	}
	if checksum.Valid {
		value := checksum.String
		attachment.Checksum = &value
	}
	return &attachment, nil
}

func (s *Service) recordAudit(ctx context.Context, organizationID, actorID, action, entityType string, entityID *string, oldValue, newValue any) error {
	if s.audit == nil {
		return nil
	}
	return s.audit.Record(ctx, organizationID, &actorID, action, entityType, entityID, oldValue, newValue, nil, nil)
}

type requestRow struct {
	ID             string
	OrganizationID string
	RequesterID    string
}

func storageDriverName(store storage.Store) string {
	switch store.(type) {
	case *storage.S3Store:
		return "s3"
	default:
		return "local"
	}
}

func storageBucketName(store storage.Store) string {
	switch s := store.(type) {
	case interface{ Bucket() string }:
		return s.Bucket()
	default:
		return "local"
	}
}

func sanitizeName(name string) string {
	name = filepath.Base(strings.TrimSpace(name))
	if name == "." || name == string(filepath.Separator) || name == "" {
		return "attachment"
	}
	var b strings.Builder
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '.', r == '-', r == '_':
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}
	out := b.String()
	if out == "" {
		return "attachment"
	}
	return out
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
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
