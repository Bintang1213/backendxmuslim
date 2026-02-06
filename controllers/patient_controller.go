package controllers

import (
	"net/http"
	"strings"
	"time"
	"strconv"
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

	if err := config.DB.First(&models.Poli{}, input.PoliID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Poli tidak valid"})
		return
	}

	// VALIDASI JAMINAN
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
			Where("nomor_jaminan = ?", *input.NomorJaminan).
			Count(&count)

		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Nomor jaminan sudah digunakan"})
			return
		}
	} else {
		input.NomorJaminan = nil
	}

	// CEK NIK
	var patient models.Patient
	err := config.DB.Where("nik = ?", input.NIK).First(&patient).Error

	// PASIEN SUDAH ADA
	if err == nil {
		if input.TipePasien == "baru" {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Pasien dengan NIK ini sudah terdaftar",
			})
			return
		}

		// PASIEN LAMA
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

	if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if input.TipePasien == "lama" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Pasien lama tidak ditemukan, silakan daftar sebagai pasien baru",
		})
		return
	}

	tgl, err := time.Parse("2006-01-02", input.TanggalLahir)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal_lahir harus yyyy-mm-dd"})
		return
	}

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
	var totalData int64

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 { page = 1 }
	if limit < 1 { limit = 10 }
	if limit > 100 { limit = 100 }
	offset := (page - 1) * limit

	tx := config.DB.Model(&models.Patient{}).Preload("Poli")

	if search := strings.TrimSpace(c.Query("search")); search != "" {
		tx = tx.Where("nama_pasien ILIKE ? OR nik LIKE ?", "%"+search+"%", search+"%")
	}

	if status := c.Query("status"); status != "" {
		tx = tx.Where("status = ?", status)
	}

	poli := c.Query("poli")
	if poli != "" && poli != "null" && poli != "undefined" {
		if poliID, err := strconv.Atoi(poli); err == nil {
			tx = tx.Where("poli_id = ?", poliID)
		} else {
			tx = tx.Joins("JOIN polis ON polis.id = patients.poli_id").
				Where("LOWER(polis.nama_poli) = ?", strings.ToLower(poli))
		}
	}

	tx.Count(&totalData)

	if err := tx.Order("created_at DESC").Limit(limit).Offset(offset).Find(&patients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pasien"})
		return
	}

	totalPages := int((totalData + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data pasien",
		"data":    patients,
		"pagination": gin.H{
			"current_page": page,
			"limit":        limit,
			"total_data":   totalData,
			"total_pages":  totalPages,
		},
	})
}
/* ================= UPDATE  ================= */

func UpdatePatient(c *gin.Context) {
    id := c.Param("id")
    var patient models.Patient

    // 1. Cari data pasien lama
    if err := config.DB.First(&patient, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Data pasien tidak ditemukan"})
        return
    }

    var input struct {
        NIK          string  `json:"nik"`
        NamaPasien   string  `json:"nama_pasien"`
        JenisKelamin string  `json:"jenis_kelamin"`
        CaraBayar    string  `json:"cara_bayar"`
        NomorJaminan *string `json:"nomor_jaminan"`
        PoliID       uint    `json:"poli_id"`
        Status       string  `json:"status"` 
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if input.NIK != "" {
        patient.NIK = input.NIK
    }
    if input.NamaPasien != "" {
        patient.NamaPasien = input.NamaPasien
    }
    if input.JenisKelamin != "" {
        patient.JenisKelamin = input.JenisKelamin
    }
    if input.Status != "" {
        patient.Status = input.Status
    }

    if input.PoliID != 0 {
        if err := config.DB.First(&models.Poli{}, input.PoliID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Poli tidak valid"})
			return
		}
        patient.PoliID = input.PoliID
    }

    if input.CaraBayar != "" {
    if input.CaraBayar == "asuransi" {
        nomor := input.NomorJaminan
        if (nomor == nil || *nomor == "") && (patient.NomorJaminan == nil || *patient.NomorJaminan == "") {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Nomor jaminan wajib diisi untuk asuransi"})
            return
        }
        if nomor != nil {
            patient.NomorJaminan = nomor
        }
    } else {
        patient.NomorJaminan = nil 
    }
    patient.CaraBayar = input.CaraBayar
}

    // 6. Simpan perubahan
    if err := config.DB.Save(&patient).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate data"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Data pasien berhasil diperbarui",
        "data":    patient,
    })
}

/* ================= DELETE ================= */

func DeletePatient(c *gin.Context) {
	id := c.Param("id")
	config.DB.Delete(&models.Patient{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Data pasien berhasil dihapus"})
}
