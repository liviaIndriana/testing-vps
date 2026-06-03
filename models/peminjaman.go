package models

import "time"


type Peminjaman struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	UserID          uint      `json:"user_id"`

	Nama            string    `json:"nama"`
	Kelas           string    `json:"kelas"`
	Tanggal         time.Time `json:"tanggal"`
	WaktuMulai      string    `json:"waktu_mulai"`
	WaktuBerakhir   string    `json:"waktu_berakhir"`
	Ruangan         string    `json:"ruangan"`
	KodeProyektor   string    `json:"kode_proyektor"`
	Keterangan      string    `json:"keterangan"`
	JenisPeminjaman string    `json:"jenis_peminjaman"`
	Status          string    `json:"status" gorm:"default:PENDING"`
}

func (Peminjaman) TableName() string {
	return "peminjaman"
}