package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/i56/framework/core/response"
)

type UploadHandler struct {
	UploadDir   string
	MaxFileSize int64 // bytes, default 32MB
}

func NewUploadHandler(dir string) *UploadHandler {
	if dir == "" {
		dir = "./uploads"
	}
	os.MkdirAll(dir, 0755)
	return &UploadHandler{UploadDir: dir, MaxFileSize: 32 << 20}
}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, nil)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, h.MaxFileSize)
	if err := r.ParseMultipartForm(h.MaxFileSize); err != nil {
		response.Error(w, fmt.Errorf("file too large (max %d MB)", h.MaxFileSize>>20))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.Error(w, fmt.Errorf("missing file field"))
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	fullPath := filepath.Join(h.UploadDir, filename)

	dst, err := os.Create(fullPath)
	if err != nil {
		response.Error(w, fmt.Errorf("failed to save file"))
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		response.Error(w, fmt.Errorf("failed to write file"))
		return
	}

	response.JSON(w, http.StatusCreated, map[string]string{
		"filename":  filename,
		"original":  header.Filename,
		"size":      fmt.Sprintf("%d", header.Size),
		"url":       "/uploads/" + filename,
	})
}
