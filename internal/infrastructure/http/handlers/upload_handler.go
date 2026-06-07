package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/alexander/go-api-hex/internal/infrastructure/storage"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	storage *storage.S3Storage
}

func NewUploadHandler(
	storage *storage.S3Storage,
) *UploadHandler {
	return &UploadHandler{
		storage: storage,
	}
}

func (h *UploadHandler) Upload(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "archivo requerido",
		})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer src.Close()

	filename := filepath.Base(file.Filename)

	url, err := h.storage.Upload(
		filename,
		src,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"filename": filename,
		"url":      url,
	})
}