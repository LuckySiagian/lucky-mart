package services

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fileHeader builds a multipart file header the way gin receives an upload.
func fileHeader(t *testing.T, name string, content []byte) *multipart.FileHeader {
	t.Helper()
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	part, err := w.CreateFormFile("file", name)
	if err != nil {
		t.Fatal(err)
	}
	part.Write(content)
	w.Close()

	req := httptest.NewRequest("POST", "/", &body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	if err := req.ParseMultipartForm(1 << 20); err != nil {
		t.Fatal(err)
	}
	return req.MultipartForm.File["file"][0]
}

// A minimal valid PNG header is enough for content sniffing.
var pngBytes = append([]byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR"), make([]byte, 64)...)

func TestSaveImageStoresPNG(t *testing.T) {
	dir := t.TempDir()
	svc := NewUploadService(dir)

	url, err := svc.SaveImage(fileHeader(t, "bukti.png", pngBytes))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(url, ImagePathPrefix) || !strings.HasSuffix(url, ".png") {
		t.Fatalf("unexpected url %q", url)
	}
	onDisk := filepath.Join(ImageDir(dir), filepath.FromSlash(strings.TrimPrefix(url, ImagePathPrefix)))
	if got, err := os.ReadFile(onDisk); err != nil || !bytes.Equal(got, pngBytes) {
		t.Fatalf("stored file mismatch: %v", err)
	}
}

func TestSaveImageRejectsNonImages(t *testing.T) {
	svc := NewUploadService(t.TempDir())
	// An HTML/script file renamed to .jpg must not be accepted.
	_, err := svc.SaveImage(fileHeader(t, "foto.jpg", []byte("<html><script>alert(1)</script></html>")))
	if !errors.Is(err, ErrImageNotAllowed) {
		t.Fatalf("want ErrImageNotAllowed, got %v", err)
	}
}

func TestIsAllowedImageURL(t *testing.T) {
	cases := map[string]bool{
		"/uploads/images/2026/10/a.jpg":    true,
		"https://example.com/a.jpg":        true,
		"/uploads/images/../../etc/passwd": false,
		"javascript:alert(1)":              false,
		"/api/v1/users":                    false,
	}
	for in, want := range cases {
		if got := IsAllowedImageURL(in); got != want {
			t.Errorf("IsAllowedImageURL(%q) = %v, want %v", in, got, want)
		}
	}
}
