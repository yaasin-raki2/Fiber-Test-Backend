package main

import (
	"ambassador/src/database"
	"ambassador/src/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	DB := database.Connect()
	database.AutoMigrate(DB)

	app := fiber.New()

	routes.Setup(app)

	app.Listen(":3000")
}
