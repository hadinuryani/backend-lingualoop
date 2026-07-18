package schedule

import (
	"errors"
	"net/http"

	"backend-lingualoop/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetConfig(c *gin.Context) {
	config, err := h.service.GetConfig(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Berhasil mengambil konfigurasi jadwal", config)
}

func (h *Handler) SaveConfig(c *gin.Context) {
	var req ScheduleConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Format input tidak valid", err.Error())
		return
	}

	config, err := h.service.SaveConfig(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Berhasil menyimpan konfigurasi jadwal", config)
}

func (h *Handler) GetAll(c *gin.Context) {
	classID := c.Query("class_id")
	academicYearID := c.Query("academic_year_id")

	var schedules []*ScheduleResponse
	var err error

	if classID != "" && academicYearID != "" {
		schedules, err = h.service.GetByClass(c.Request.Context(), classID, academicYearID)
	} else {
		schedules, err = h.service.GetAll(c.Request.Context())
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Berhasil mengambil data jadwal", schedules)
}

func (h *Handler) Create(c *gin.Context) {
	var req ScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Format input tidak valid", err.Error())
		return
	}

	sch, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrClassClash) || errors.Is(err, ErrTeacherClash) {
			statusCode = http.StatusConflict
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusCreated, "Berhasil membuat jadwal", sch)
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req ScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Format input tidak valid", err.Error())
		return
	}

	sch, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrScheduleNotFound) {
			statusCode = http.StatusNotFound
		} else if errors.Is(err, ErrClassClash) || errors.Is(err, ErrTeacherClash) {
			statusCode = http.StatusConflict
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Berhasil mengubah jadwal", sch)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrScheduleNotFound) {
			statusCode = http.StatusNotFound
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Berhasil menghapus jadwal", nil)
}

func (h *Handler) DeleteByClass(c *gin.Context) {
	classID := c.Param("class_id")
	academicYearID := c.Query("academic_year_id")
	
	if academicYearID == "" {
		response.Error(c, http.StatusBadRequest, "academic_year_id wajib diisi", nil)
		return
	}

	if err := h.service.DeleteByClass(c.Request.Context(), classID, academicYearID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Berhasil menghapus jadwal kelas", nil)
}
