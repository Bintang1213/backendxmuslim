package controllers

import (
	"net/http"
	"strings"
	"time"

	"tugas12/config"
	"tugas12/models"

	"github.com/gin-gonic/gin"
)

/* ================= CREATE ================= */

type CreatePatientInput struct {
	NIK          string  `json:"nik" binding:"required,len=16"`
	NamaPasien   string  `json:"nama_pasien" binding:"required"`
	JenisKelamin string  `json:"jenis_kelamin" binding:"required,oneof=L P"`
	TanggalLahir string  `json:"tanggal_lahir" binding:"required"` // yyyy-mm-dd
	TipePasien   string  `json:"tipe_pasien" binding:"required,oneof=baru lama"`
	CaraBayar    string  `json:"cara_bayar" binding:"required,oneof=umum asuransi"`
	NomorJaminan *string `json:"nomor_jaminan"`
	PoliID       uint    `json:"poli_id" binding:"required"`
	Status       string  `json:"status"` // proses / selesai (opsional)
}

func CreatePatient(c *gin.Context) {
	var input CreatePatientInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse tanggal lahir (biar gak 2000-01-01T00:00:00Z)
	tglLahir, err := time.Parse("2006-01-02", input.TanggalLahir)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal_lahir harus yyyy-mm-dd"})
		return
	}

	// Default status
	status := input.Status
	if status == "" {
		status = "proses"
	}

	if status != "proses" && status != "selesai" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status harus proses atau selesai"})
		return
	}

	// Validasi poli (dropdown dari DB)
	var poli models.Poli
	if err := config.DB.First(&poli, input.PoliID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Poli tidak valid"})
		return
	}

	// Validasi asuransi
	if input.CaraBayar == "asuransi" {
		if input.NomorJaminan == nil || len(*input.NomorJaminan) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nomor jaminan wajib diisi"})
			return
		}
	} else {
		input.NomorJaminan = nil
	}

	patient := models.Patient{
		NIK:          input.NIK,
		NamaPasien:   input.NamaPasien,
		JenisKelamin: input.JenisKelamin,
		TanggalLahir: tglLahir,
		TipePasien:   input.TipePasien,
		CaraBayar:    input.CaraBayar,
		NomorJaminan: input.NomorJaminan,
		PoliID:       input.PoliID,
		Status:       status,
	}

	if err := config.DB.Create(&patient).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pasien berhasil ditambahkan",
		"data":    patient,
	})
}

/* ================= READ ================= */

func GetPatients(c *gin.Context) {
	var patients []models.Patient

	search := c.Query("search")
	status := c.Query("status")
	poli := c.Query("poli")

	db := config.DB.Preload("Poli")

	if search != "" {
		db = db.Where("LOWER(nama_pasien) LIKE ?", "%"+strings.ToLower(search)+"%")
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}
	if poli != "" {
		db = db.Where("poli_id = ?", poli)
	}

	db.Order("created_at DESC").Find(&patients)
	c.JSON(http.StatusOK, patients)
}

/* ================= UPDATE ================= */

func UpdatePatient(c *gin.Context) {
	id := c.Param("id")
	var patient models.Patient

	if err := config.DB.First(&patient, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pasien tidak ditemukan"})
		return
	}

	var input CreatePatientInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tglLahir, err := time.Parse("2006-01-02", input.TanggalLahir)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal_lahir harus yyyy-mm-dd"})
		return
	}

	patient.NamaPasien = input.NamaPasien
	patient.JenisKelamin = input.JenisKelamin
	patient.TanggalLahir = tglLahir
	patient.TipePasien = input.TipePasien
	patient.CaraBayar = input.CaraBayar
	patient.NomorJaminan = input.NomorJaminan
	patient.PoliID = input.PoliID
	patient.Status = input.Status

	config.DB.Save(&patient)
	c.JSON(http.StatusOK, gin.H{"message": "Pasien berhasil diperbarui"})
}

/* ================= DELETE ================= */

func DeletePatient(c *gin.Context) {
	id := c.Param("id")

	if err := config.DB.Delete(&models.Patient{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pasien berhasil dihapus"})
}
