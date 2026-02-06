package main

import (
	"os"
	"time" 
	"tugas12/config"
	"tugas12/models"
	"tugas12/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config.ConnectDB()

	//  MIGRATION
	config.DB.AutoMigrate(
		&models.User{},
		&models.Poli{},
		&models.Patient{},
	)

	seedPoli()

	r := gin.Default()

	//  CORS Configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"https://cedrick-unlunated-gwyn.ngrok-free.app", 
		},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin", 
			"Content-Type", 
			"Authorization", 
			"Accept", 
			"ngrok-skip-browser-warning", 
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, 
	}))

	routes.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

func seedPoli() {
	polis := []string{
		"Jantung",
		"Mata",
		"THT",
		"Anak",
		"Penyakit Dalam",
	}

	for _, nama := range polis {
		config.DB.FirstOrCreate(&models.Poli{}, models.Poli{NamaPoli: nama})
	}
}