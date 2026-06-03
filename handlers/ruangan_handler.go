package handlers

import (
	"backend-codegirls/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RuanganHandler struct {
	DB *gorm.DB
}
// GET ALL RUANGAN
func (h *RuanganHandler) GetRuangan(c *fiber.Ctx) error {
	var data []models.Ruangan

	if err := h.DB.Find(&data).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Server error",
		})
	}

	return c.JSON(fiber.Map{
		"data": data,
	})
}

// UPDATE RUANGAN
func (h *RuanganHandler) UpdateRuangan(c *fiber.Ctx) error {
	id := c.Params("id")

	var body models.Ruangan

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Input tidak valid",
		})
	}

	// validasi
	if body.Ruangan == "" || body.Fasilitas == "" || body.Kapasitas <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Semua field wajib diisi",
		})
	}

	var ruangan models.Ruangan

	if err := h.DB.First(&ruangan, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Ruangan tidak ditemukan",
		})
	}

	ruangan.Ruangan = body.Ruangan
	ruangan.Fasilitas = body.Fasilitas
	ruangan.Kapasitas = body.Kapasitas

	h.DB.Save(&ruangan)

	return c.JSON(fiber.Map{
		"message": "Ruangan berhasil diupdate",
	})
}


// DELETE RUANGAN
func (h *RuanganHandler) DeleteRuangan(c *fiber.Ctx) error {
	id := c.Params("id")

	var ruangan models.Ruangan

	if err := h.DB.First(&ruangan, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Ruangan tidak ditemukan",
		})
	}

	h.DB.Delete(&ruangan)

	return c.JSON(fiber.Map{
		"message": "Ruangan berhasil dihapus",
	})
}