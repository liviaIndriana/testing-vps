package handlers

import (
	"backend-codegirls/models"
	"backend-codegirls/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Register(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user models.User


		if err := c.BodyParser(&user); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid input",
			})
		}

		// validasi
		if user.Nama == "" || user.NPMNIDN == "" || user.Email == "" || user.Password == "" || user.Role == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Semua field wajib diisi",
			})
		}

		// cek NPM/NIDN sudah ada atau belum
		var existing models.User
		if err := db.Where("npm_nidn = ?", user.NPMNIDN).First(&existing).Error; err == nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "NPM/NIDN sudah digunakan",
			})
		}
		if err := db.Where("email = ?", user.Email).First(&existing).Error; err == nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Email sudah digunakan",
			})
		}
		// HASH PASSWORD
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to hash password",
			})
		}

		user.Password = hashedPassword

		// simpan ke database
		if err := db.Create(&user).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Gagal menyimpan user",
			})
		}

		return c.JSON(fiber.Map{
			"message": "User created",
		})
	}
}

func Login(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type LoginInput struct {
			NPMNIDN string `json:"npm_nidn"`
			Password string `json:"password"`
		}

		var input LoginInput
		var user models.User

		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid input",
			})
		}

		// validasi
		if input.NPMNIDN == "" || input.Password == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "NPM/NIDN dan password wajib diisi",
			})
		}

		// cari user berdasarkan NPM/NIDN
		if err := db.Where("npm_nidn = ?", input.NPMNIDN).First(&user).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": "User tidak ditemukan",
			})
		}

		// cek password
		if !utils.CheckPassword(user.Password, input.Password) {
			return c.Status(401).JSON(fiber.Map{
				"error": "Password salah",
			})
		}

		// generate token
		token, err := utils.GenerateToken(user.ID, user.NPMNIDN, string(user.Role))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Gagal generate token",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Login success",
			"token":   token,
			"user": fiber.Map{
				"id":       user.ID,
				"nama":     user.Nama,
				"npm_nidn": user.NPMNIDN,
				"role":     user.Role,
			},
		})
	}
}