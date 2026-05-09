package storage

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

type S3Store struct {
	endpoint       string
	bucket         string
	accessKey      string
	secretKey      string
	region         string
	forcePathStyle bool
	client         *http.Client
	bucketEnsured  bool
}

func (s *S3Store) Bucket() string {
	return s.bucket
}

func NewS3Store(endpoint, bucket, accessKey, secretKey, region string, forcePathStyle bool) *S3Store {
	return &S3Store{
		endpoint:       strings.TrimRight(endpoint, "/"),
		bucket:         bucket,
		accessKey:      accessKey,
		secretKey:      secretKey,
		region:         region,
		forcePathStyle: forcePathStyle,
		client:         &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *S3Store) EnsureBucket(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, s.bucketURL(""), http.NoBody)
	if err != nil {
		return err
	}
	s.signRequest(req, []byte{})
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 && resp.StatusCode != http.StatusConflict {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ensure bucket failed: %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}
	return nil
}

func (s *S3Store) Put(ctx context.Context, key, contentType string, data []byte) error {
	if err := s.EnsureBucket(ctx); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, s.objectURL(key), bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	s.signRequest(req, data)
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("put object failed: %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}
	return nil
}

func (s *S3Store) Get(ctx context.Context, key string) ([]byte, error) {
	if err := s.EnsureBucket(ctx); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.objectURL(key), http.NoBody)
	if err != nil {
		return nil, err
	}
	s.signRequest(req, []byte{})
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get object failed: %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}
	return io.ReadAll(resp.Body)
}

func (s *S3Store) bucketURL(key string) string {
	parsed, _ := url.Parse(s.endpoint)
	if s.forcePathStyle {
		parsed.Path = path.Join("/", s.bucket, key)
		return parsed.String()
	}
	host := parsed.Host
	parsed.Host = s.bucket + "." + host
	parsed.Path = path.Join("/", key)
	return parsed.String()
}

func (s *S3Store) objectURL(key string) string {
	return s.bucketURL(key)
}

func (s *S3Store) signRequest(req *http.Request, payload []byte) {
	now := time.Now().UTC()
	amzDate := now.Format("20060102T150405Z")
	shortDate := now.Format("20060102")
	payloadHash := sha256Hex(payload)

	req.Header.Set("x-amz-date", amzDate)
	req.Header.Set("x-amz-content-sha256", payloadHash)
	req.Host = req.URL.Host
	if req.Host == "" {
		req.Host = req.URL.Hostname()
	}

	canonicalURI := req.URL.EscapedPath()
	if canonicalURI == "" {
		canonicalURI = "/"
	}
	canonicalHeaders := "host:" + strings.ToLower(req.Host) + "\n" +
		"x-amz-content-sha256:" + payloadHash + "\n" +
		"x-amz-date:" + amzDate + "\n"
	signedHeaders := "host;x-amz-content-sha256;x-amz-date"
	canonicalRequest := strings.Join([]string{
		req.Method,
		canonicalURI,
		req.URL.RawQuery,
		canonicalHeaders,
		signedHeaders,
		payloadHash,
	}, "\n")

	scope := shortDate + "/" + s.region + "/s3/aws4_request"
	stringToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256",
		amzDate,
		scope,
		sha256Hex([]byte(canonicalRequest)),
	}, "\n")

	signingKey := s.signingKey(shortDate)
	signature := hex.EncodeToString(hmacSHA256(signingKey, []byte(stringToSign)))
	req.Header.Set("Authorization", fmt.Sprintf(
		"AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		s.accessKey, scope, signedHeaders, signature,
	))
}

func (s *S3Store) signingKey(date string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+s.secretKey), []byte(date))
	kRegion := hmacSHA256(kDate, []byte(s.region))
	kService := hmacSHA256(kRegion, []byte("s3"))
	return hmacSHA256(kService, []byte("aws4_request"))
}

func hmacSHA256(key, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write(data)
	return mac.Sum(nil)
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
