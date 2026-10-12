package handler

import (
	"net/http"

	"backend-sembako/helpers"
	"backend-sembako/internal/services"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	uploadService services.UploadService
}

func NewUploadHandler(uploadService services.UploadService) *UploadHandler {
	return &UploadHandler{uploadService: uploadService}
}

// UploadImage stores one photo sent as multipart field "file" and returns its URL.
func (h *UploadHandler) UploadImage(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, services.MaxImageSize+1<<20)
	file, err := c.FormFile("file")
	if err != nil {
		helpers.SendBadRequest(c, "Pilih foto terlebih dahulu (maksimal 5 MB).", nil)
		return
	}

	url, err := h.uploadService.SaveImage(file)
	if err != nil {
		sendServiceError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusCreated, "Foto berhasil diunggah", gin.H{"url": url})
}
