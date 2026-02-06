package models

import "time"

type Patient struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    NIK          string    `gorm:"size:16;not null" json:"nik"`
    NamaPasien   string    `gorm:"not null" json:"nama_pasien"`
    JenisKelamin string    `gorm:"size:15;not null" json:"jenis_kelamin"`
    TanggalLahir time.Time `gorm:"type:date;not null" json:"tanggal_lahir"`
    TipePasien   string    `gorm:"not null" json:"tipe_pasien"`
    CaraBayar    string    `gorm:"not null" json:"cara_bayar"`
    NomorJaminan *string   `gorm:"size:20" json:"nomor_jaminan"`
    PoliID       uint      `gorm:"not null" json:"poli_id"`
    Poli         Poli      `gorm:"foreignKey:PoliID" json:"poli"` 
    Status       string    `gorm:"not null;default:'proses'" json:"status"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}