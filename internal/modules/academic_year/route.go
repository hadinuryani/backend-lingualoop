package academic_year

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

func RegisterRoute(router *gin.RouterGroup, db *sql.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	templateService := NewTemplateService(repo)
	handler := NewHandler(service, templateService)

	ay := router.Group("/academic-years")
	{
		ay.GET("", handler.GetAll)
		ay.GET("/template", handler.GetTemplate)
		ay.GET("/:id", handler.GetByID)
		ay.POST("", handler.Create)
		ay.PUT("/:id", handler.Update)
		ay.DELETE("/:id", handler.Delete)

		ay.PUT("/:id/activate", handler.Activate)
		ay.PUT("/:id/semester-status", handler.UpdateSemesterStatus)
		ay.PUT("/:id/close-semester", handler.CloseSemester)

		// Kenaikan kelas ditaruh di sini nanti
		// ay.POST("/:id/finalize-promotion", handler.FinalizePromotion)
	}
}
