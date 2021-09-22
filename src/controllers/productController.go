package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Products(c *fiber.Ctx) error {
	var products []models.Product

	database.DB.Find(&products)

	return c.JSON(products)
}

func CreacteProducts(c *fiber.Ctx) error {
	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return err
	}

	database.DB.Create(&product)

	return c.JSON(product)
}

func GetProduct(c *fiber.Ctx) error {
	var product models.Product

	id, _ := strconv.Atoi(c.Params("id"))

	product.Id = uint(id)

	database.DB.Find(&product)

	return c.JSON(product)
}

func UpdateProduct(c *fiber.Ctx) error {
	var product models.Product

	id, _ := strconv.Atoi(c.Params("id"))

	product.Id = uint(id)

	if err := c.BodyParser(&product); err != nil {
		return err
	}

	database.DB.Model(&product).Updates(&product).Find(&product)

	return c.JSON(product)
}

func DeleteProduct(c *fiber.Ctx) error {
	var product models.Product

	id, _ := strconv.Atoi(c.Params("id"))

	product.Id = uint(id)

	database.DB.Delete(&product)

	return nil
}

func ProductsFrontend(c *fiber.Ctx) error {
	var products []models.Product
	var ctx = context.Background()

	result, err := database.Cache.Get(ctx, "products_frontend").Result()

	if err != nil {
		database.DB.Find(&products)

		bytes, err := json.Marshal(products)

		if err != nil {
			panic(err)
		}

		if err = database.Cache.Set(ctx, "products_frontend", bytes, 30*time.Minute).Err(); err != nil {
			panic(err)
		}
	} else {
		json.Unmarshal([]byte(result), &products)
	}

	return c.JSON(&products)
}

func ProductsBackend(c *fiber.Ctx) error {
	var products []models.Product
	var ctx = context.Background()

	result, err := database.Cache.Get(ctx, "products_backend").Result()

	if err != nil {
		database.DB.Find(&products)

		bytes, _ := json.Marshal(products)

		database.Cache.Set(ctx, "products_backend", bytes, 30*time.Minute)
	} else {
		json.Unmarshal([]byte(result), &products)
	}

	var searchedProducts []models.Product

	if s := c.Query("s"); s != "" {
		lowerS := strings.ToLower(s)

		for _, product := range products {
			lowerT := strings.ToLower(product.Title)
			lowerD := strings.ToLower(product.Description)

			if strings.Contains(lowerT, lowerS) || strings.Contains(lowerD, lowerS) {
				searchedProducts = append(searchedProducts, product)
			}
		}
	} else {
		return c.JSON(&products)
	}

	return c.JSON(&searchedProducts)
}
