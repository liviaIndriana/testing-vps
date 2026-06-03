package routes

import (
	"backend-codegirls/handlers"
	"backend-codegirls/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api")

	// auth
	api.Post("/register", handlers.Register(db))
	api.Post("/login", handlers.Login(db))

	// ADMIN GROUP
	admin := api.Group("/admin", middleware.AdminOnly())

	// dashboard
	admin.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome Admin",
		})
	})

	// HISTORY
	historyHandler := &handlers.HistoryHandler{DB: db}
	app.Get("/history", middleware.AdminOnly(), historyHandler.GetHistory)
	app.Patch("/history/:id", middleware.AdminOnly(), historyHandler.UpdateStatus)
	admin.Put("/history/:id/approve", historyHandler.Approve)
	admin.Put("/history/:id/reject", historyHandler.Reject)
	app.Post("/peminjaman",middleware.AuthMiddleware(),historyHandler.CreatePeminjaman,)
	app.Get("/history/me",middleware.AuthMiddleware(),historyHandler.GetMyHistory,)


	// RUANGAN
	ruanganHandler := &handlers.RuanganHandler{DB: db}
	api.Get("/ruangan", middleware.AdminOnly(), ruanganHandler.GetRuangan)
	api.Put("/ruangan/:id", middleware.AdminOnly(), ruanganHandler.UpdateRuangan)
	api.Delete("/ruangan/:id", middleware.AdminOnly(), ruanganHandler.DeleteRuangan)
	api.Get("/user-ruangan",ruanganHandler.GetRuangan,)

	// JADWAL
	jadwalHandler := &handlers.JadwalHandler{DB: db}
	api.Post("/jadwal", middleware.AdminOnly(), jadwalHandler.CreateJadwal)
	api.Get("/jadwal/terjadwal", jadwalHandler.GetTerjadwal)
	api.Get("/jadwal/tidak-terjadwal", jadwalHandler.GetTidakTerjadwal)
	

	// MANAJEMEN USER
	manajemenHandler := &handlers.ManajemenHandler{DB: db}
	api.Get("/manajemen",middleware.AdminOnly(),manajemenHandler.GetUsers,)
	api.Delete("/manajemen/:id_user",middleware.AdminOnly(),manajemenHandler.DeleteUser,)
	api.Patch("/manajemen/reset-password/:id_user",middleware.AdminOnly(),manajemenHandler.ResetPassword,)
}
