package handlers

import (
	"backend-codegirls/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type JadwalHandler struct {
	DB *gorm.DB
}

func (h *JadwalHandler) CreateJadwal(c *fiber.Ctx) error {

	var body struct {
		Kelas         string `json:"kelas"`
		Tanggal       string `json:"tanggal"`
		WaktuMulai    string `json:"waktu_mulai"`
		WaktuBerakhir string `json:"waktu_berakhir"`
		Ruangan       string `json:"ruangan"`
		Jenis       string `json:"jenis"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input tidak valid"})
	}

	// VALIDASI WAJIB
	if body.Kelas == "" || body.Tanggal == "" || body.WaktuMulai == "" || body.WaktuBerakhir == "" || body.Ruangan == "" {
		return c.Status(400).JSON(fiber.Map{"message": "Semua field wajib diisi"})
	}

	// VALIDASI TANGGAL
	tanggalParsed, err := time.Parse("2006-01-02", body.Tanggal)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Format tanggal salah"})
	}

	// VALIDASI WAKTU
	if body.WaktuBerakhir <= body.WaktuMulai {
		return c.Status(400).JSON(fiber.Map{"message": "Waktu tidak valid"})
	}

	// CEK BENTROK DENGAN JADWAL LAIN
	var existingJadwal models.Jadwal
	err = h.DB.Where(
		"ruangan = ? AND tanggal = ? AND ((waktu_mulai < ? AND waktu_berakhir > ?) OR (waktu_mulai < ? AND waktu_berakhir > ?))",
		body.Ruangan,
		tanggalParsed,
		body.WaktuBerakhir, body.WaktuMulai,
		body.WaktuMulai, body.WaktuBerakhir,
	).First(&existingJadwal).Error

	if err == nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Jadwal bentrok dengan jadwal lain",
		})
	}

	//  CEK BENTROK DENGAN PEMINJAMAN
	var existingPinjam models.Peminjaman
	err = h.DB.Where(
		"ruangan = ? AND DATE(tanggal) = ? AND status = ? AND NOT (waktu_berakhir <= ? OR waktu_mulai >= ?)",
		body.Ruangan,
		body.Tanggal,
		"APPROVED",
		body.WaktuMulai,
		body.WaktuBerakhir,
	).First(&existingPinjam).Error

		// SIMPAN
	jenis := body.Jenis

	if jenis != "terjadwal" && jenis != "tidak terjadwal" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Jenis tidak valid",
		})
	}

	jadwal := models.Jadwal{
		Kelas:         body.Kelas,
		Tanggal:       tanggalParsed,
		WaktuMulai:    body.WaktuMulai,
		WaktuBerakhir: body.WaktuBerakhir,
		Ruangan:       body.Ruangan,
		Jenis:         jenis,
	}

	if err := h.DB.Create(&jadwal).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal menyimpan"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Jadwal berhasil ditambahkan",
	})
}

func (h *JadwalHandler) GetTerjadwal(c *fiber.Ctx) error {
	var jadwal []models.Jadwal

	hari := c.Query("hari")
	hariMap := map[string]string{
		"Senin":  "Monday",
		"Selasa": "Tuesday",
		"Rabu":   "Wednesday",
		"Kamis":  "Thursday",
		"Jumat":  "Friday",
		"Sabtu":  "Saturday",
		"Minggu": "Sunday",
	}
	query := h.DB.Where("jenis = ?", "terjadwal")
	if hari != "" {

		hariDB, ok := hariMap[hari]

		if ok {
			query = query.Where(
				"DAYNAME(tanggal) = ?",
				hariDB,
			)
		}
	}

	err := query.Find(&jadwal).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal mengambil data",
		})
	}

	var result []fiber.Map

	for _, j := range jadwal {
		result = append(result, fiber.Map{
			"id_jadwal":      j.ID,
			"kelas":          j.Kelas,
			"ruangan":        j.Ruangan,
			"tanggal":        j.Tanggal.Format("2006-01-02"),
			"waktu_mulai":    j.WaktuMulai,
			"waktu_berakhir": j.WaktuBerakhir,
		})
	}

	return c.JSON(fiber.Map{
		"terjadwal": result,
	})
}

func (h *JadwalHandler) GetTidakTerjadwal(c *fiber.Ctx) error {

	var jadwal []models.Jadwal

	hari := c.Query("hari")

	hariMap := map[string]string{
		"Senin":  "Monday",
		"Selasa": "Tuesday",
		"Rabu":   "Wednesday",
		"Kamis":  "Thursday",
		"Jumat":  "Friday",
		"Sabtu":  "Saturday",
		"Minggu": "Sunday",
	}

	query := h.DB.Where(
		"jenis = ?", "tidak terjadwal",
	)

	if hari != "" {

		hariDB, ok := hariMap[hari]

		if ok {
			query = query.Where(
				"DAYNAME(tanggal) = ?",
				hariDB,
			)
		}
	}

	err := query.Find(&jadwal).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal mengambil data",
		})
	}

	var result []fiber.Map

	for _, p := range jadwal {

		result = append(result, fiber.Map{
			"id_jadwal":      p.ID,
			"kelas":          p.Kelas,
			"ruangan":        p.Ruangan,
			"tanggal":        p.Tanggal.Format("2006-01-02"),
			"waktu_mulai":    p.WaktuMulai,
			"waktu_berakhir": p.WaktuBerakhir,
		})
	}

	return c.JSON(fiber.Map{
		"tidak_terjadwal": result,
	})
}
