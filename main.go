package main

import (
	"ambassador/src/database"

	"github.com/gofiber/fiber/v2"
)

func main() {
	DB := database.Connect()
	database.AutoMigrate(DB)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello 1337 !")
	})

	app.Listen(":3000")
}
