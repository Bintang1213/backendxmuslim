package dto

type CreatePatientRequest struct {
	NIK          string  `json:"nik" binding:"required,len=16"`
	NamaPasien   string  `json:"nama_pasien" binding:"required"`
	JenisKelamin string  `json:"jenis_kelamin" binding:"required,oneof=L P"`
	TanggalLahir string  `json:"tanggal_lahir" binding:"required"` 
	TipePasien   string  `json:"tipe_pasien" binding:"required,oneof=baru lama"`
	CaraBayar    string  `json:"cara_bayar" binding:"required,oneof=umum asuransi"`
	NomorJaminan *string `json:"nomor_jaminan"`
	PoliID       uint    `json:"poli_id" binding:"required"`
}
