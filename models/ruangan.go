package models

type Ruangan struct {
	ID         uint   `gorm:"primaryKey" json:"id_ruangan"`
	Ruangan    string `json:"nama_ruangan"`
	Fasilitas  string `json:"fasilitas"`
	Kapasitas  int    `json:"kapasitas"`
}
func (Ruangan) TableName() string {
	return "ruangan"
}