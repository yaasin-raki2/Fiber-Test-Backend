package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"

	"github.com/gofiber/fiber/v2"
)

func Orders(c *fiber.Ctx) error {
	var orders []models.Order

	database.DB.Preload("OrderItems").Find(&orders)

	for i := range orders {
		orders[i].Name = orders[i].FullName()
		orders[i].Total = orders[i].GetTotal()
	}

	return c.JSON(orders)
}
