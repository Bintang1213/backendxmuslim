package main

import (
	"os"
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

	// 🔹 MIGRATION
	config.DB.AutoMigrate(
		&models.User{},
		&models.Poli{},
		&models.Patient{},
	)

	seedPoli()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
		},
		AllowHeaders: []string{"Authorization", "Content-Type"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
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
