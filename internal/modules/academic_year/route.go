package academic_year

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	
	"backend-lingualoop/internal/modules/class"
	"backend-lingualoop/internal/modules/student"
)

func RegisterRoute(router *gin.RouterGroup, db *sql.DB) {
	repo := NewRepository(db)
	classRepo := class.NewRepository(db)
	studentRepo := student.NewRepository(db)
	
	service := NewService(repo, classRepo, studentRepo, db)
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

		ay.POST("/:id/finalize-promotion", handler.FinalizePromotion)
	}
}
