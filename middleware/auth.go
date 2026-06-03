package middleware


import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var SECRET_KEY = []byte("secret123")

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token tidak ditemukan",
			})
		}

		tokenString := strings.Split(authHeader, " ")
		if len(tokenString) != 2 {
			return c.Status(401).JSON(fiber.Map{
				"error": "Format token salah",
			})
		}

		token, err := jwt.Parse(tokenString[1], func(token *jwt.Token) (interface{}, error) {
			return SECRET_KEY, nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token tidak valid",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token tidak valid",
			})
		}

		// cek role
		if claims["role"] != "admin" {
			return c.Status(403).JSON(fiber.Map{
				"error": "Akses khusus admin",
			})
		}

		return c.Next()
	}
}

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token tidak ditemukan",
			})
		}

		tokenString := strings.Split(authHeader, " ")

		if len(tokenString) != 2 {
			return c.Status(401).JSON(fiber.Map{
				"error": "Format token salah",
			})
		}

		token, err := jwt.Parse(tokenString[1], func(token *jwt.Token) (interface{}, error) {
			return SECRET_KEY, nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token tidak valid",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token tidak valid",
			})
		}

		// simpan user login ke locals
		c.Locals("user_id", uint(claims["user_id"].(float64)))
		c.Locals("role", claims["role"])

		return c.Next()
	}
}

