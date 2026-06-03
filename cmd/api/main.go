package main

import (
	"log"
	"os"

	"backend-codegirls/config"
	"backend-codegirls/models"
	"backend-codegirls/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// koneksi database
	db := config.NewDatabase()

	// migrate tabel user
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("❌ Gagal migrate tabel User:", err)
	}
	log.Println("✅ Berhasil migrate tabel User")

	// migrate tabel peminjaman
	if err := db.AutoMigrate(&models.Peminjaman{}); err != nil {
		log.Fatal("❌ Gagal migrate tabel Peminjaman:", err)
	}
	log.Println("✅ Berhasil migrate tabel Peminjaman")

	// migrate tabel ruangan
	if err := db.AutoMigrate(&models.Ruangan{}); err != nil {
		log.Fatal("❌ Gagal migrate tabel Ruangan:", err)
	}
	log.Println("✅ Berhasil migrate tabel Ruangan")

	// migrate tabel jadwal
	if err := db.AutoMigrate(&models.Jadwal{}); err != nil {
		log.Fatal("❌ Gagal migrate tabel Jadwal:", err)
	}
	log.Println("✅ Berhasil migrate tabel Jadwal")

	// init app
	app := fiber.New()

	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS,PATCH",
	}))

	// routes
	routes.SetupRoutes(app, db)

	// port
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3001"
	}

	log.Println("Server running on port", port)
	app.Listen(":" + port)
}
