package routes

import (
	"ambassador/src/controllers"
	"ambassador/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("api")

	admin := api.Group("admin")
	admin.Post("/register", controllers.Register)
	admin.Post("/login", controllers.Login)

	adminAuthenticated := admin.Use(middlewares.IsAuthenticated)
	adminAuthenticated.Get("/user", controllers.User)
	adminAuthenticated.Post("/logout", controllers.Logout)
	adminAuthenticated.Put("/user/info", controllers.UpdateInfo)
	adminAuthenticated.Put("/user/password", controllers.UpdatePassword)
	adminAuthenticated.Get("/ambassadors", controllers.Ambassadors)
	adminAuthenticated.Get("/products", controllers.Products)
	adminAuthenticated.Post("/products", controllers.CreacteProducts)
	adminAuthenticated.Get("/products/:id", controllers.GetProduct)
	adminAuthenticated.Put("/products/:id", controllers.UpdateProduct)
	adminAuthenticated.Delete("/products/:id", controllers.DeleteProduct)
	adminAuthenticated.Get("/user/:id/links", controllers.Link)
	adminAuthenticated.Get("/orders", controllers.Orders)

	ambassador := api.Group("ambassador")
	ambassador.Post("/register", controllers.Register)
	ambassador.Post("/login", controllers.Login)
	ambassador.Get("/products/frontend", controllers.ProductsFrontend)

	ambassadorAuthenticated := ambassador.Use(middlewares.IsAuthenticated)
	ambassadorAuthenticated.Get("/user", controllers.User)
	ambassadorAuthenticated.Post("/logout", controllers.Logout)
	ambassadorAuthenticated.Put("/user/info", controllers.UpdateInfo)
	ambassadorAuthenticated.Put("/user/password", controllers.UpdatePassword)
}
