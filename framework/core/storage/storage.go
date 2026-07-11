// Package storage provides unified file storage abstraction.
package storage

import (
	"context"
	"io"
	"time"
)

// PutOptions configures file upload behavior.
type PutOptions struct {
	ContentType string
	PublicRead  bool
	Metadata    map[string]string
}

// Storage is the abstract file storage interface.
type Storage interface {
	Put(ctx context.Context, path string, data io.Reader, opts *PutOptions) error
	Get(ctx context.Context, path string) (io.ReadCloser, error)
	Delete(ctx context.Context, path string) error
	URL(ctx context.Context, path string, expiry time.Duration) (string, error)
	Exists(ctx context.Context, path string) (bool, error)
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
