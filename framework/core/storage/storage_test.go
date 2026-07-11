package storage

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
	"time"
)

// inMemStorage implements Storage for testing.
type inMemStorage struct {
	files map[string][]byte
}

func newInMemStorage() *inMemStorage {
	return &inMemStorage{files: make(map[string][]byte)}
}

func (s *inMemStorage) Put(ctx context.Context, path string, data io.Reader, opts *PutOptions) error {
	b, err := io.ReadAll(data)
	if err != nil {
		return err
	}
	s.files[path] = b
	return nil
}

func (s *inMemStorage) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	b, ok := s.files[path]
	if !ok {
		return nil, &storageError{"file not found"}
	}
	return io.NopCloser(bytes.NewReader(b)), nil
}

func (s *inMemStorage) Delete(ctx context.Context, path string) error {
	delete(s.files, path)
	return nil
}

func (s *inMemStorage) URL(ctx context.Context, path string, expiry time.Duration) (string, error) {
	return "http://example.com/" + path, nil
}

func (s *inMemStorage) Exists(ctx context.Context, path string) (bool, error) {
	_, ok := s.files[path]
	return ok, nil
}

type storageError struct{ msg string }

func (e *storageError) Error() string { return e.msg }

func TestStorage_PutAndGet(t *testing.T) {
	s := newInMemStorage()
	ctx := context.Background()

	data := strings.NewReader("hello world")
	err := s.Put(ctx, "file.txt", data, &PutOptions{ContentType: "text/plain"})
	if err != nil {
		t.Fatalf("Put: %v", err)
	}

	rc, err := s.Get(ctx, "file.txt")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	defer rc.Close()

	b, _ := io.ReadAll(rc)
	if string(b) != "hello world" {
		t.Errorf("expected 'hello world', got %q", b)
	}
}

func TestStorage_Exists(t *testing.T) {
	s := newInMemStorage()
	ctx := context.Background()

	ok, err := s.Exists(ctx, "nonexistent.txt")
	if err != nil {
		t.Fatalf("Exists: %v", err)
	}
	if ok {
		t.Error("expected false for nonexistent file")
	}

	s.Put(ctx, "exists.txt", strings.NewReader("data"), nil)

	ok, err = s.Exists(ctx, "exists.txt")
	if err != nil {
		t.Fatalf("Exists: %v", err)
	}
	if !ok {
		t.Error("expected true for existing file")
	}
}

func TestStorage_Delete(t *testing.T) {
	s := newInMemStorage()
	ctx := context.Background()

	s.Put(ctx, "file.txt", strings.NewReader("data"), nil)
	s.Delete(ctx, "file.txt")

	ok, _ := s.Exists(ctx, "file.txt")
	if ok {
		t.Error("expected file to be deleted")
	}
}

func TestStorage_URL(t *testing.T) {
	s := newInMemStorage()
	ctx := context.Background()

	url, err := s.URL(ctx, "path/to/file.pdf", time.Hour)
	if err != nil {
		t.Fatalf("URL: %v", err)
	}
	if url == "" {
		t.Error("expected non-empty URL")
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
