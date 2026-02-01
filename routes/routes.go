package routes

import (
	"tugas12/controllers"
	"tugas12/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/login", controllers.Login)
		api.POST("/register", controllers.Register)

		api.GET("/polis", middleware.JWTAuth(), controllers.GetPolis)

		api.GET("/patients", middleware.JWTAuth(), controllers.GetPatients)
		api.POST("/patients", middleware.JWTAuth(), controllers.CreatePatient)
		api.PUT("/patients/:id", middleware.JWTAuth(), controllers.UpdatePatient)
		api.DELETE("/patients/:id", middleware.JWTAuth(), controllers.DeletePatient)
	}
}
