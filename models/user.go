package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Nama     string `json:"nama"`
	NPMNIDN  string `gorm:"unique;column:npm_nidn" json:"npm_nidn"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}