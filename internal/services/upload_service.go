package services

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// MaxImageSize is the largest photo accepted (phone photos are usually 1–4 MB).
const MaxImageSize = 5 << 20

// ImagePathPrefix is the public URL path uploaded images are served from.
const ImagePathPrefix = "/uploads/images/"

var (
	ErrImageTooLarge   = errors.New("ukuran foto maksimal 5 MB")
	ErrImageNotAllowed = errors.New("file harus berupa foto JPG, PNG, atau WebP")
	ErrImageEmpty      = errors.New("pilih foto terlebih dahulu")

	errInvalidProductImage = errors.New("foto produk tidak valid; unggah foto atau gunakan link https://")
)

// Only real image formats are accepted; the type is detected from the file
// content, never trusted from the file name or browser.
var imageExtensions = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

type UploadService interface {
	SaveImage(file *multipart.FileHeader) (string, error)
}

type uploadService struct {
	imageDir string
}

func NewUploadService(uploadDir string) UploadService {
	return &uploadService{imageDir: filepath.Join(uploadDir, "images")}
}

// ImageDir is where images are stored on disk, for serving them back.
func ImageDir(uploadDir string) string {
	return filepath.Join(uploadDir, "images")
}

// SaveImage stores an uploaded photo under images/YYYY/MM with a random name
// and returns its public path, e.g. /uploads/images/2026/10/<uuid>.jpg.
func (s *uploadService) SaveImage(fh *multipart.FileHeader) (string, error) {
	if fh == nil || fh.Size == 0 {
		return "", ErrImageEmpty
	}
	if fh.Size > MaxImageSize {
		return "", ErrImageTooLarge
	}

	src, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	head := make([]byte, 512)
	n, err := io.ReadFull(src, head)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, io.EOF) {
		return "", err
	}
	ext, ok := imageExtensions[http.DetectContentType(head[:n])]
	if !ok {
		return "", ErrImageNotAllowed
	}

	sub := time.Now().Format("2006/01")
	dir := filepath.Join(s.imageDir, filepath.FromSlash(sub))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("gagal menyiapkan folder unggahan: %w", err)
	}

	name := uuid.NewString() + ext
	dst, err := os.OpenFile(filepath.Join(dir, name), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return "", err
	}
	// Copy what was already sniffed, then the rest, capped at the size limit.
	written, err := io.Copy(dst, io.LimitReader(io.MultiReader(bytes.NewReader(head[:n]), src), MaxImageSize+1))
	closeErr := dst.Close()
	if err == nil {
		err = closeErr
	}
	if err == nil && written > MaxImageSize {
		err = ErrImageTooLarge
	}
	if err != nil {
		_ = os.Remove(filepath.Join(dir, name))
		return "", err
	}

	return ImagePathPrefix + sub + "/" + name, nil
}

// IsAllowedImageURL accepts a photo uploaded to this server or an external
// http(s) link (kept for older data and pasted links).
func IsAllowedImageURL(u string) bool {
	u = strings.TrimSpace(u)
	if strings.HasPrefix(u, ImagePathPrefix) {
		return !strings.Contains(u, "..")
	}
	return strings.HasPrefix(u, "https://") || strings.HasPrefix(u, "http://")
}
