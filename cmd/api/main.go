package main

import (
	"log"
	"os"

	"backend-codegirls/config"
	"backend-codegirls/routes"
	"backend-codegirls/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// koneksi database
	db := config.NewDatabase()

	// migrate tabel user
	db.AutoMigrate(&models.User{})

	// migrate tabel peminjaman
	db.AutoMigrate(&models.Peminjaman{})

	// migrate tabel ruangan
	db.AutoMigrate(&models.Ruangan{})

	// mrigate tabel jadwal
	db.AutoMigrate(&models.Jadwal{})

	// init app
	app := fiber.New()

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