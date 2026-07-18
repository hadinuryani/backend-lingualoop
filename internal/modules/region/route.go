package region

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoute(router *gin.RouterGroup, db *sql.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	regionRoute := router.Group("/region")
	{
		// Provinces
		regionRoute.GET("/provinces", handler.GetAllProvinces)
		regionRoute.GET("/provinces/:id", handler.GetProvinceByID)

		// Cities (nested under province)
		regionRoute.GET("/provinces/:provinceId/cities", handler.GetCitiesByProvinceID)
		regionRoute.GET("/cities/:id", handler.GetCityByID)

		// Districts (nested under city)
		regionRoute.GET("/cities/:cityId/districts", handler.GetDistrictsByCityID)
		regionRoute.GET("/districts/:id", handler.GetDistrictByID)

		// Subdistricts (nested under district)
		regionRoute.GET("/districts/:districtId/subdistricts", handler.GetSubdistrictsByDistrictID)
		regionRoute.GET("/subdistricts/:id", handler.GetSubdistrictByID)

		// Postal Codes
		regionRoute.GET("/subdistricts/:subdistrictId/postal-codes", handler.GetPostalCodesBySubdistrictID)
		regionRoute.GET("/postal-codes/search", handler.SearchPostalCode)
	}
}
