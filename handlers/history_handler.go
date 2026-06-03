package handlers

import (
	"backend-codegirls/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HistoryHandler struct {
	DB *gorm.DB
}

//
// =======================
// CREATE PEMINJAMAN (INSERT)
// =======================
//
func (h *HistoryHandler) CreatePeminjaman(c *fiber.Ctx) error {
	var body struct {
		Nama            string `json:"nama"`
		Kelas           string `json:"kelas"`
		Tanggal         string `json:"tanggal"`
		WaktuMulai      string `json:"waktu_mulai"`
		WaktuBerakhir   string `json:"waktu_berakhir"`
		KodeProyektor   string `json:"kode_proyektor"`
		Keterangan      string `json:"keterangan"`
		JenisPeminjaman string `json:"jenis_peminjaman"`
		Ruangan         string `json:"ruangan"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Request tidak valid",
		})
	}

	// parse tanggal
	tgl, err := time.Parse("2006-01-02", body.Tanggal)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Format tanggal harus YYYY-MM-DD",
		})
	}
	userID := c.Locals("user_id").(uint)

	data := models.Peminjaman{
		UserID:          userID,
		Nama:            body.Nama,
		Kelas:           body.Kelas,
		Tanggal:         tgl,
		WaktuMulai:      body.WaktuMulai,
		WaktuBerakhir:   body.WaktuBerakhir,
		KodeProyektor:   body.KodeProyektor,
		Keterangan:      body.Keterangan,
		JenisPeminjaman: body.JenisPeminjaman,
		Status:          "PENDING",
		Ruangan:         body.Ruangan,
	}

	if err := h.DB.Create(&data).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal membuat peminjaman",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Peminjaman berhasil dibuat",
		"data":    data,
	})
}

//
// =======================
// GET HISTORY
// =======================
//
func (h *HistoryHandler) GetHistory(c *fiber.Ctx) error {
	var data []models.Peminjaman

	if err := h.DB.Find(&data).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Server error",
		})
	}

	var result []fiber.Map

	for _, item := range data {
		result = append(result, fiber.Map{
			"id":               item.ID,
			"nama":             item.Nama,
			"kelas":            item.Kelas,
			"tanggal":          item.Tanggal.Format("2006-01-02"),
			"waktu_mulai":      item.WaktuMulai,
			"waktu_berakhir":   item.WaktuBerakhir,
			"kode_proyektor":   item.KodeProyektor,
			"keterangan":       item.Keterangan,
			"jenis_peminjaman": item.JenisPeminjaman,
			"ruangan":          item.Ruangan,
			"status":           item.Status,
		})
	}

	return c.JSON(result)
}

//
// =======================
// UPDATE STATUS (GENERAL)
// =======================
//
func (h *HistoryHandler) UpdateStatus(c *fiber.Ctx) error {
	id := c.Params("id")

	var body struct {
		Status string `json:"status"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Request tidak valid",
		})
	}

	var data models.Peminjaman
	if err := h.DB.First(&data, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Data tidak ditemukan",
		})
	}

	if body.Status != "APPROVED" && body.Status != "REJECTED" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Status tidak valid",
		})
	}

	data.Status = body.Status

	if err := h.DB.Save(&data).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal update status",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Status berhasil diupdate",
	})
}

//
// =======================
// APPROVE
// =======================
//
func (h *HistoryHandler) Approve(c *fiber.Ctx) error {
	id := c.Params("id")

	var data models.Peminjaman
	if err := h.DB.First(&data, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Data tidak ditemukan",
		})
	}

	data.Status = "APPROVED"
	h.DB.Save(&data)

	return c.JSON(fiber.Map{
		"message": "Peminjaman disetujui",
	})
}

//
// =======================
// REJECT
// =======================
//
func (h *HistoryHandler) Reject(c *fiber.Ctx) error {
	id := c.Params("id")

	var data models.Peminjaman
	if err := h.DB.First(&data, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Data tidak ditemukan",
		})
	}

	data.Status = "REJECTED"
	h.DB.Save(&data)

	return c.JSON(fiber.Map{
		"message": "Peminjaman ditolak",
	})
}

//
// =======================
// HISTORY USER
// =======================
//

func (h *HistoryHandler) GetMyHistory(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var data []models.Peminjaman

	if err := h.DB.
		Where("user_id = ?", userID).
		Find(&data).Error; err != nil {

		return c.Status(500).JSON(fiber.Map{
			"message": "Server error",
		})
	}

	var result []fiber.Map

	for _, item := range data {
		result = append(result, fiber.Map{
			"id":               item.ID,
			"nama":             item.Nama,
			"kelas":            item.Kelas,
			"tanggal":          item.Tanggal.Format("2006-01-02"),
			"waktu_mulai":      item.WaktuMulai,
			"waktu_berakhir":   item.WaktuBerakhir,
			"kode_proyektor":   item.KodeProyektor,
			"keterangan":       item.Keterangan,
			"jenis_peminjaman": item.JenisPeminjaman,
			"ruangan":          item.Ruangan,
			"status":           item.Status,
		})
	}

	return c.JSON(result)
}

