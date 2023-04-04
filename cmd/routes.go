package main

import (
	"template/handlers"
	"template/middleware"

	"github.com/dgraph-io/ristretto"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App, cache *ristretto.Cache) {
	// Product routes
	app.Get("/product/:id", handlers.GetProductHandler(cache))
	app.Get("/", handlers.GetAllProductsHandler(cache))
	app.Get("/api/products", handlers.GetFilteredProductsHandler(cache))
	app.Post("/api/product", handlers.CreateProduct)
	app.Get("/tmp", handlers.Tmp)
	// User routes
	app.Post("/api/register", handlers.Register)
	app.Post("/api/login", handlers.Login)

	// Protected route example
	app.Use("/api/whoami", middleware.Protected())
	app.Get("/api/whoami", handlers.WhoAmI)
}
