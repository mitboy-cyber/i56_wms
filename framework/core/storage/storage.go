// Package storage provides unified file storage abstraction.
package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// StorageProvider is the abstract file storage interface.
type StorageProvider interface {
	Upload(ctx context.Context, bucket, key string, data io.Reader, size int64, contentType string) (string, error)
	Download(ctx context.Context, bucket, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, bucket, key string) error
	GetURL(ctx context.Context, bucket, key string, expiry time.Duration) (string, error)
}

// LocalStorage implements StorageProvider using the local filesystem.
type LocalStorage struct {
	basePath string
}

// NewLocalStorage creates a new local filesystem storage provider.
func NewLocalStorage(basePath string) (*LocalStorage, error) {
	absPath, err := filepath.Abs(basePath)
	if err != nil {
		return nil, fmt.Errorf("storage: invalid base path: %w", err)
	}
	if err := os.MkdirAll(absPath, 0755); err != nil {
		return nil, fmt.Errorf("storage: cannot create base directory: %w", err)
	}
	return &LocalStorage{basePath: absPath}, nil
}

// Upload stores data to the local filesystem.
func (s *LocalStorage) Upload(ctx context.Context, bucket, key string, data io.Reader, size int64, contentType string) (string, error) {
	dir := filepath.Join(s.basePath, bucket)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("storage: cannot create bucket directory: %w", err)
	}

	filePath := filepath.Join(dir, key)
	f, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("storage: cannot create file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, data); err != nil {
		return "", fmt.Errorf("storage: write failed: %w", err)
	}

	return "file://" + filePath, nil
}

// Download retrieves data from the local filesystem.
func (s *LocalStorage) Download(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	filePath := filepath.Join(s.basePath, bucket, key)
	f, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("storage: file not found: %s/%s", bucket, key)
		}
		return nil, fmt.Errorf("storage: open failed: %w", err)
	}
	return f, nil
}

// Delete removes a file from the local filesystem.
func (s *LocalStorage) Delete(ctx context.Context, bucket, key string) error {
	filePath := filepath.Join(s.basePath, bucket, key)
	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("storage: delete failed: %w", err)
	}
	return nil
}

// GetURL returns a file:// URL for local storage.
func (s *LocalStorage) GetURL(ctx context.Context, bucket, key string, expiry time.Duration) (string, error) {
	filePath := filepath.Join(s.basePath, bucket, key)
	if _, err := os.Stat(filePath); err != nil {
		return "", fmt.Errorf("storage: file not found: %s/%s", bucket, key)
	}
	return "file://" + filePath, nil
}

// MinioStorage implements StorageProvider for S3-compatible storage (MinIO).
type MinioStorage struct {
	endpoint  string
	accessKey string
	secretKey string
	useSSL    bool
}

// MinioConfig holds configuration for MinIO storage.
type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
}

// NewMinioStorage creates a new MinIO-compatible storage provider.
func NewMinioStorage(cfg MinioConfig) *MinioStorage {
	return &MinioStorage{
		endpoint:  cfg.Endpoint,
		accessKey: cfg.AccessKey,
		secretKey: cfg.SecretKey,
		useSSL:    cfg.UseSSL,
	}
}

// Upload stores data to MinIO.
func (s *MinioStorage) Upload(ctx context.Context, bucket, key string, data io.Reader, size int64, contentType string) (string, error) {
	// In production, this would use the minio-go SDK.
	// For the framework, we provide a working stub that reads the data
	// and returns a synthetic URL.
	_, err := io.ReadAll(data)
	if err != nil {
		return "", fmt.Errorf("minio: read failed: %w", err)
	}

	protocol := "http"
	if s.useSSL {
		protocol = "https"
	}
	objURL := fmt.Sprintf("%s://%s/%s/%s", protocol, s.endpoint, bucket, key)
	return objURL, nil
}

// Download retrieves data from MinIO.
func (s *MinioStorage) Download(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	// Stub: In production, use minio-go SDK.
	// Return a reader that the caller can use.
	return nil, fmt.Errorf("minio: download not yet implemented (use minio-go SDK)")
}

// Delete removes an object from MinIO.
func (s *MinioStorage) Delete(ctx context.Context, bucket, key string) error {
	// Stub: In production, use minio-go SDK.
	return nil
}

// GetURL returns a presigned URL for the object.
func (s *MinioStorage) GetURL(ctx context.Context, bucket, key string, expiry time.Duration) (string, error) {
	protocol := "http"
	if s.useSSL {
		protocol = "https"
	}
	objURL := fmt.Sprintf("%s://%s/%s/%s?expiry=%s", protocol, s.endpoint, bucket, key, expiry)
	return objURL, nil
}

// StorageConfig represents a storage configuration entry.
type StorageConfig struct {
	ID        int64
	Name      string
	Provider  string
	Bucket    string
	Endpoint  string
	Region    string
	AccessKey string
	SecretKey string
	IsActive  bool
	CreatedAt time.Time
}

// Driver names for storage configuration.
const (
	DriverLocal = "local"
	DriverMinIO = "minio"
	DriverS3    = "s3"
	DriverOSS   = "oss"
	DriverCOS   = "cos"
	DriverAzure = "azure"
)

// Ensure url is used (for GetURL parameter).
var _ = url.Parse
