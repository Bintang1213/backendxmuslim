package controllers

import (
	"net/http"
	"tugas12/config"
	"tugas12/models"

	"github.com/gin-gonic/gin"
)

func GetPolis(c *gin.Context) {
	var polis []models.Poli
	config.DB.Find(&polis)
	c.JSON(http.StatusOK, polis)
}
