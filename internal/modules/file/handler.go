package file

import (
	"context"
	"net/http"
	"time"

	"backend-lingualoop/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// UploadImage handles image upload specifically
func (h *Handler) UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "File tidak ditemukan dalam request", nil)
		return
	}
	defer file.Close()

	resourceType := c.PostForm("resource_type")
	if resourceType == "" {
		response.Error(c, http.StatusBadRequest, "resource_type wajib diisi", nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// Optionally get user ID if authenticated
	var uploadedBy *string
	// userID, exists := c.Get("user_id")
	// if exists {
	// 	idStr := userID.(string)
	// 	uploadedBy = &idStr
	// }

	res, err := h.service.UploadImage(ctx, file, header, resourceType, uploadedBy)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == ErrInvalidFileType || err == ErrFileTooLarge {
			statusCode = http.StatusBadRequest
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Berhasil mengunggah file", res)
}
