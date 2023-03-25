package main

import (
	"template/handlers"

	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {

	app.Get("/product/:id", handlers.GetProduct)
	app.Get("/", handlers.GetAllProducts)

	app.Get("/api/products", handlers.GetFilteredProducts)
	app.Post("/api/product", handlers.CreateProduct)
}
