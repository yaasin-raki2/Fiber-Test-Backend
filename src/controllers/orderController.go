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

type CreateOrderRequest struct {
	Code      string
	FirstName string
	LastName  string
	Email     string
	Address   string
	Country   string
	City      string
	Zip       string
	Products  []map[string]int
}

func CreateOrder(c *fiber.Ctx) error {
	var request CreateOrderRequest

	if err := c.BodyParser(&request); err != nil {
		return err
	}

	link := models.Link{
		Code: request.Code,
	}

	database.DB.Preload("User").First(&link)

	if link.Id == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid link!",
		})
	}

	order := models.Order{
		UserId:          link.UserId,
		Code:            request.Code,
		AmbassadorEmail: link.User.Email,
		FirstName:       request.FirstName,
		LastName:        request.LastName,
		Email:           request.Email,
		Address:         request.Address,
		Country:         request.Country,
		City:            request.City,
		Zip:             request.Zip,
	}

	tx := database.DB.Begin()

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	for _, requestProduct := range request.Products {
		product := models.Product{}
		product.Id = uint(requestProduct["product_id"])
		database.DB.First(&product)

		total := product.Price * float64(requestProduct["quantity"])

		item := models.OrderItem{
			OrderId:           order.Id,
			ProductTitle:      product.Title,
			Price:             product.Price,
			Quantity:          uint(requestProduct["quantity"]),
			AmbassadorRevenue: 0.1 * total,
			AdminRevenue:      0.9 * total,
		}

		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	}

	tx.Commit()

	return c.JSON(order)
}
