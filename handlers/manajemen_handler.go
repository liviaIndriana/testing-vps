package handlers

import (
	"backend-codegirls/models"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

type ManajemenHandler struct {
	DB *gorm.DB
}

// GET ALL USER
func (h *ManajemenHandler) GetUsers(c *fiber.Ctx) error {
	var users []models.User

	if err := h.DB.Find(&users).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal mengambil data user",
		})
	}

	var result []fiber.Map

	for _, u := range users {
		result = append(result, fiber.Map{
			"id_user":  u.ID,
			"nama":     u.Nama,
			"npm_nidn": u.NPMNIDN,
			"email":    u.Email,
			"role":     u.Role,
		})
	}

	return c.JSON(fiber.Map{
		"data": result,
	})
}

// DELETE USER
func (h *ManajemenHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id_user")

	var user models.User

	// cek user berdasarkan id
	if err := h.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "User tidak ditemukan",
		})
	}

	// admin tidak boleh dihapus
	role := strings.ToLower(strings.TrimSpace(user.Role))

	if role == "admin" {
		return c.Status(403).JSON(fiber.Map{
			"message": "Admin tidak bisa dihapus",
		})
	}

	// hapus user
	if err := h.DB.Delete(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal menghapus user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User berhasil dihapus",
	})
}

func (h *ManajemenHandler) ResetPassword(c *fiber.Ctx) error {

	id := c.Params("id_user")

	var user models.User

	// cari user
	if err := h.DB.
		Where("id = ?", id).
		First(&user).Error; err != nil {

		return c.Status(404).JSON(fiber.Map{
			"message": "User tidak ditemukan",
		})
	}

	// password baru = npm/nidn
	newPassword := user.NPMNIDN

	// hash password
	hashedPassword, err :=
		bcrypt.GenerateFromPassword(
			[]byte(newPassword),
			bcrypt.DefaultCost,
		)

	if err != nil {

		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal hash password",
		})
	}

	user.Password = string(hashedPassword)

	if err := h.DB.Save(&user).Error; err != nil {

		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal reset password",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Password berhasil direset",
	})
}