package storage

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"
)

func TestLocalStorage_PutAndGet(t *testing.T) {
	s, err := NewLocalStorage(t.TempDir())
	if err != nil {
		t.Fatalf("NewLocalStorage: %v", err)
	}
	ctx := context.Background()

	data := strings.NewReader("hello world")
	_, err = s.Upload(ctx, "test-bucket", "file.txt", data, int64(len("hello world")), "text/plain")
	if err != nil {
		t.Fatalf("Upload: %v", err)
	}

	rc, err := s.Download(ctx, "test-bucket", "file.txt")
	if err != nil {
		t.Fatalf("Download: %v", err)
	}
	defer rc.Close()

	b, _ := io.ReadAll(rc)
	if string(b) != "hello world" {
		t.Errorf("expected 'hello world', got %q", b)
	}
}

func TestLocalStorage_Delete(t *testing.T) {
	s, _ := NewLocalStorage(t.TempDir())
	ctx := context.Background()

	s.Upload(ctx, "bucket", "file.txt", strings.NewReader("data"), 4, "text/plain")

	err := s.Delete(ctx, "bucket", "file.txt")
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}

	_, err = s.Download(ctx, "bucket", "file.txt")
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestLocalStorage_GetURL(t *testing.T) {
	s, _ := NewLocalStorage(t.TempDir())
	ctx := context.Background()

	s.Upload(ctx, "bucket", "file.txt", strings.NewReader("data"), 4, "text/plain")

	url, err := s.GetURL(ctx, "bucket", "file.txt", time.Hour)
	if err != nil {
		t.Fatalf("GetURL: %v", err)
	}
	if url == "" {
		t.Error("expected non-empty URL")
	}
}

func TestLocalStorage_DownloadNotFound(t *testing.T) {
	s, _ := NewLocalStorage(t.TempDir())
	ctx := context.Background()

	_, err := s.Download(ctx, "bucket", "nonexistent.txt")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestDriverConstants(t *testing.T) {
	if DriverLocal != "local" {
		t.Errorf("expected 'local', got %q", DriverLocal)
	}
	if DriverMinIO != "minio" {
		t.Errorf("expected 'minio', got %q", DriverMinIO)
	}
	if DriverS3 != "s3" {
		t.Errorf("expected 's3', got %q", DriverS3)
	}
}

func TestMinioStorage_Upload(t *testing.T) {
	s := NewMinioStorage(MinioConfig{
		Endpoint:  "minio.example.com:9000",
		AccessKey: "minioadmin",
		SecretKey: "minioadmin",
		UseSSL:    false,
	})
	ctx := context.Background()

	data := strings.NewReader("test data")
	url, err := s.Upload(ctx, "bucket", "key", data, 9, "text/plain")
	if err != nil {
		t.Fatalf("Upload: %v", err)
	}
	if url == "" {
		t.Error("expected non-empty URL")
	}
}

func TestMinioStorage_GetURL(t *testing.T) {
	s := NewMinioStorage(MinioConfig{
		Endpoint:  "s3.example.com",
		AccessKey: "key",
		SecretKey: "secret",
		UseSSL:    true,
	})
	ctx := context.Background()

	url, err := s.GetURL(ctx, "bucket", "key", time.Hour)
	if err != nil {
		t.Fatalf("GetURL: %v", err)
	}
	if !strings.Contains(url, "https://") {
		t.Errorf("expected HTTPS URL, got %q", url)
	}
}
