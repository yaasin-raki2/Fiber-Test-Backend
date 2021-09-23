package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"context"
	"encoding/json"
	"sort"
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

	go database.ClearCache("products_frontend", "products_backend")

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

	go database.ClearCache("products_frontend", "products_backend")

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
		searchedProducts = products
	}

	if sortParam := c.Query("sort"); sortParam != "" {
		sortLower := strings.ToLower(sortParam)
		if sortLower == "asc" {
			sort.Slice(searchedProducts, func(i, j int) bool {
				return searchedProducts[i].Price < searchedProducts[j].Price
			})
		} else if sortLower == "desc" {
			sort.Slice(searchedProducts, func(i, j int) bool {
				return searchedProducts[i].Price > searchedProducts[j].Price
			})
		}
	}

	total := len(searchedProducts)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perpage := 9

	if page < 1 {
		page = 1
	} else if page > total/perpage+1 {
		page = total/perpage + 1
	}

	var data []models.Product

	if total <= page*perpage && total >= (page-1)*perpage {
		data = searchedProducts[(page-1)*perpage : total]
	} else if total >= page*perpage {
		data = searchedProducts[(page-1)*perpage : page*perpage]
	} else {
		data = []models.Product{}
	}

	return c.JSON(fiber.Map{
		"data":      data,
		"total":     total,
		"page":      page,
		"last_page": total/perpage + 1,
	})
}
