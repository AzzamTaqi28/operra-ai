package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

type LocalStore struct {
	baseDir string
}

func NewLocalStore(baseDir string) *LocalStore {
	if baseDir == "" {
		baseDir = filepath.Join(os.TempDir(), "operra-storage")
	}
	return &LocalStore{baseDir: baseDir}
}

func (s *LocalStore) EnsureBucket(ctx context.Context) error {
	return os.MkdirAll(s.baseDir, 0o755)
}

func (s *LocalStore) Put(ctx context.Context, key, contentType string, data []byte) error {
	path := filepath.Join(s.baseDir, key)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func (s *LocalStore) Get(ctx context.Context, key string) ([]byte, error) {
	path := filepath.Join(s.baseDir, key)
	return os.ReadFile(path)
}

func (s *LocalStore) Bucket() string {
	return "local"
}

func sanitizeKey(parts ...string) string {
	key := ""
	for i, part := range parts {
		if i > 0 {
			key += "/"
		}
		key += clean(part)
	}
	return key
}

func clean(value string) string {
	out := make([]rune, 0, len(value))
	for _, r := range value {
		switch {
		case r >= 'a' && r <= 'z':
			out = append(out, r)
		case r >= 'A' && r <= 'Z':
			out = append(out, r)
		case r >= '0' && r <= '9':
			out = append(out, r)
		case r == '-', r == '_', r == '.', r == '/':
			out = append(out, r)
		default:
			out = append(out, '-')
		}
	}
	return string(out)
}

func localFileKey(orgID, requestID, attachmentID, fileName string) string {
	return sanitizeKey("attachments", orgID, requestID, fmt.Sprintf("%s_%s", attachmentID, fileName))
}
