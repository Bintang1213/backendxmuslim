package models

import "time"

type Patient struct {
	ID            uint      `gorm:"primaryKey"`
	NIK           string    `gorm:"size:16;not null"`
	NamaPasien    string    `gorm:"not null"`
	JenisKelamin  string    `gorm:"size:1;not null"`
	TanggalLahir  time.Time `gorm:"type:date;not null"`
	TipePasien    string    `gorm:"not null"`
	CaraBayar     string    `gorm:"not null"`
	NomorJaminan  *string   `gorm:"size:20"`
	PoliID        uint      `gorm:"not null"`
	Poli          Poli      `gorm:"foreignKey:PoliID"`
	Status        string    `gorm:"not null;default:'proses'"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
