package models

import "time"

type Jadwal struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Kelas         string    `json:"kelas"`
	Tanggal       time.Time `json:"tanggal"`
	WaktuMulai    string    `json:"waktu_mulai"`
	WaktuBerakhir string    `json:"waktu_berakhir"`
	Ruangan       string    `json:"ruangan"`
	Jenis 		  string    `gorm:"column:jenis" json:"jenis"` 
}

func (Jadwal) TableName() string {
	return "jadwal"
}