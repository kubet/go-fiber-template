package main

import (
	"log"

	"template/database"

	"github.com/dgraph-io/ristretto"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		log.Fatalf("Error creating cache: %v", err)
	}

	database.ConnectDb()
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // Set the desired compression level
	}))

	app.Static("/static", "/static", fiber.Static{
		Compress:      true,  // Enable on-the-fly compression for static files
		ByteRange:     true,  // Enable byte-range requests (for large files like videos)
		MaxAge:        86400, // Set cache-control header to 24 hours (86400 seconds) for better performance
		CacheDuration: 86400, // Set the expires header to 24 hours for better performance
	})

	setupRoutes(app, cache)

	app.Listen(":3080")
}
