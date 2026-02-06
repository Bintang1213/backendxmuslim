package controllers

import (
	"net/http"
	"tugas12/config"
	"tugas12/models"

	"github.com/gin-gonic/gin"
)

func GetPolis(c *gin.Context) {
	var polis []models.Poli
	if err := config.DB.Find(&polis).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data poli"})
		return
	}
	c.JSON(http.StatusOK, polis)
}