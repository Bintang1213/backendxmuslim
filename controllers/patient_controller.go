package controllers

import (
	"net/http"
	"strings"
	"time"

	"tugas12/config"
	"tugas12/dto"
	"tugas12/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/* ================= CREATE / DAFTAR ================= */

func CreatePatient(c *gin.Context) {
	var input dto.CreatePatientRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	/* ===== VALIDASI POLI ===== */
	if err := config.DB.First(&models.Poli{}, input.PoliID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Poli tidak valid"})
		return
	}

	/* ===== VALIDASI CARA BAYAR ===== */
	if input.CaraBayar == "asuransi" {
		if input.NomorJaminan == nil || len(*input.NomorJaminan) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nomor jaminan wajib diisi"})
			return
		}
		if len(*input.NomorJaminan) > 20 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nomor jaminan maksimal 20 karakter"})
			return
		}

		// cek nomor jaminan unik
		var count int64
		config.DB.Model(&models.Patient{}).
			Where("nomor_jaminan = ?", *input.NomorJaminan).
			Count(&count)

		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Nomor jaminan sudah digunakan"})
			return
		}
	} else {
		input.NomorJaminan = nil
	}

	/* ===== CEK PASIEN BERDASARKAN NIK ===== */
	var patient models.Patient
	err := config.DB.Where("nik = ?", input.NIK).First(&patient).Error

	/* ===== PASIEN SUDAH ADA ===== */
	if err == nil {

		// ❌ Pasien baru tapi NIK sudah ada
		if input.TipePasien == "baru" {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Pasien dengan NIK ini sudah terdaftar",
			})
			return
		}

		// ✅ Pasien lama → hanya update data tertentu
		patient.CaraBayar = input.CaraBayar
		patient.NomorJaminan = input.NomorJaminan
		patient.PoliID = input.PoliID
		patient.Status = "proses"

		config.DB.Save(&patient)

		c.JSON(http.StatusOK, gin.H{
			"message": "Pasien lama berhasil didaftarkan ulang",
			"data":    patient,
		})
		return
	}

	/* ===== ERROR DATABASE ===== */
	if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	/* ===== PASIEN BELUM ADA ===== */
	if input.TipePasien == "lama" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Pasien lama tidak ditemukan, silakan daftar sebagai pasien baru",
		})
		return
	}

	/* ===== PARSE TANGGAL LAHIR ===== */
	tgl, err := time.Parse("2006-01-02", input.TanggalLahir)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal_lahir harus yyyy-mm-dd"})
		return
	}

	/* ===== CREATE PASIEN BARU ===== */
	newPatient := models.Patient{
		NIK:          input.NIK,
		NamaPasien:   input.NamaPasien,
		JenisKelamin: input.JenisKelamin,
		TanggalLahir: tgl,
		TipePasien:   "baru",
		CaraBayar:    input.CaraBayar,
		NomorJaminan: input.NomorJaminan,
		PoliID:       input.PoliID,
		Status:       "proses",
	}

	config.DB.Create(&newPatient)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pasien baru berhasil ditambahkan",
		"data":    newPatient,
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Data pasien tidak ditemukan"})
		return
	}

	var input struct {
		CaraBayar    string  `json:"cara_bayar" binding:"required,oneof=umum asuransi"`
		NomorJaminan *string `json:"nomor_jaminan"`
		PoliID       uint    `json:"poli_id" binding:"required"`
		Status       string  `json:"status" binding:"required,oneof=proses selesai"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	/* ===== VALIDASI POLI ===== */
	if err := config.DB.First(&models.Poli{}, input.PoliID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Poli tidak valid"})
		return
	}

	/* ===== VALIDASI JAMINAN ===== */
	if input.CaraBayar == "asuransi" {
		if input.NomorJaminan == nil || len(*input.NomorJaminan) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nomor jaminan wajib diisi"})
			return
		}
		if len(*input.NomorJaminan) > 20 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nomor jaminan maksimal 20 karakter"})
			return
		}

		var count int64
		config.DB.Model(&models.Patient{}).
			Where("nomor_jaminan = ? AND id != ?", *input.NomorJaminan, patient.ID).
			Count(&count)

		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Nomor jaminan sudah digunakan"})
			return
		}
	} else {
		input.NomorJaminan = nil
	}

	/* ===== UPDATE DATA ===== */
	patient.CaraBayar = input.CaraBayar
	patient.NomorJaminan = input.NomorJaminan
	patient.PoliID = input.PoliID
	patient.Status = input.Status

	config.DB.Save(&patient)

	c.JSON(http.StatusOK, gin.H{"message": "Data pasien berhasil diperbarui"})
}

/* ================= DELETE ================= */

func DeletePatient(c *gin.Context) {
	id := c.Param("id")

	config.DB.Delete(&models.Patient{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Data pasien berhasil dihapus"})
}
