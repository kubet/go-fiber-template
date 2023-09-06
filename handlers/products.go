package handlers

import (
	"fmt"
	"template/database"
	"template/models"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/gofiber/fiber/v2"
)

type PageData struct {
	Title    string
	Heading  string
	Products []models.Product
}

type SinglePageData struct {
	Title               string
	Heading             string
	Product             models.Product
	RecommendedProducts []models.Product // new field to hold recommended products
}

func GetProduct(c *fiber.Ctx) error {
	productId := c.Params("id")

	product := models.Product{}
	database.DB.Db.Where("id = ?", productId).First(&product)

	data := SinglePageData{
		Title:   product.Name,
		Heading: product.Name,
		Product: product,
	}

	return c.Render("p2", data)
}

func Tmp(c *fiber.Ctx) error {
	return c.Render("tmp", PageData{})
}

func GetAllProducts(c *fiber.Ctx) error {
	products := []models.Product{}
	database.DB.Db.Where("price >= 50 AND price % 2 = 0 AND (name like '%0' or name like '%5')").Limit(100).Find(&products)

	data := PageData{
		Title:    "List of Products",
		Heading:  "Products",
		Products: products,
	}
	return c.Render("index", data)
}
func GetFilteredProducts(c *fiber.Ctx) error {
	searchQuery := c.Query("search")
	products := []models.Product{}

	if searchQuery != "" {
		database.DB.Db.Where("name LIKE ?", "%"+searchQuery).Limit(100).Find(&products)
	} else {
		database.DB.Db.Where("price >= 50 AND price % 2 = 0 AND (name like '%0' or name like '%5')").Limit(100).Find(&products)
	}

	return c.JSON(products)
}

func CreateProduct(c *fiber.Ctx) error {
	product := new(models.Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	database.DB.Db.Create(&product)

	return c.Status(200).JSON(product)
}

func productCacheKey(id string) string {
	return fmt.Sprintf("product_%s", id)
}

func allProductsCacheKey() string {
	return "all_products"
}
func filteredProductsCacheKey(searchQuery string) string {
	return fmt.Sprintf("filtered_products_%s", searchQuery)
}
func GetFilteredProductsHandler(cache *ristretto.Cache) fiber.Handler {
	return func(c *fiber.Ctx) error {
		searchQuery := c.Query("search")
		cacheKey := filteredProductsCacheKey(searchQuery)

		cachedProducts, found := cache.Get(cacheKey)

		var products []models.Product

		if found {
			var ok bool
			products, ok = cachedProducts.([]models.Product)
			if !ok {
				// Handle type assertion error here, e.g., log it
				fmt.Printf("Type assertion error: %v\n", cachedProducts)
			}
		} else {
			if searchQuery != "" {
				database.DB.Db.Where("name LIKE ?", "%"+searchQuery).Limit(100).Find(&products)
			} else {
				database.DB.Db.Where("price >= 50 AND price % 2 = 0 AND (name like '%0' or name like '%5')").Limit(100).Find(&products)
			}
			cache.SetWithTTL(cacheKey, products, 1, time.Minute*10) // Cache for 10 minutes
		}

		return c.JSON(products)
	}
}
func getRandomProducts(excludeID uint, limit int) []models.Product {
	var products []models.Product
	database.DB.Db.Where("id <> ?", excludeID).Order("RANDOM()").Limit(limit).Find(&products)
	return products
}

func GetProductHandler(cache *ristretto.Cache) fiber.Handler {
	return func(c *fiber.Ctx) error {
		productId := c.Params("id")
		cacheKey := productCacheKey(productId)

		cachedProduct, found := cache.Get(cacheKey)

		var product models.Product

		if found {
			var ok bool
			product, ok = cachedProduct.(models.Product)
			if !ok {
				fmt.Printf("Type assertion error: %v\n", cachedProduct)
			}
		} else {
			database.DB.Db.Where("id = ?", productId).First(&product)
			cache.SetWithTTL(cacheKey, product, 1, time.Minute*10) // Cache for 10 minutes
		}

		recommendedProducts := getRandomProducts(product.Id, 5) // Get 5 random recommended products excluding the current product

		data := SinglePageData{
			Title:               product.Name,
			Heading:             product.Name,
			Product:             product,
			RecommendedProducts: recommendedProducts, // include recommended products in the data
		}

		return c.Render("p2", data) // Adjusted template name to 'product'
	}
}

func GetAllProductsHandler(cache *ristretto.Cache) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cacheKey := allProductsCacheKey()

		cachedProducts, found := cache.Get(cacheKey)

		var products []models.Product

		if found {
			var ok bool
			products, ok = cachedProducts.([]models.Product)
			if !ok {
				// Handle type assertion error here, e.g., log it
				fmt.Printf("Type assertion error: %v\n", cachedProducts)
			}
		} else {
			database.DB.Db.Where("price >= 50 AND price % 2 = 0 AND (name like '%0' or name like '%5')").Limit(100).Find(&products)
			cache.SetWithTTL(cacheKey, products, 1, time.Minute*10) // Cache for 10 minutes
		}

		data := PageData{
			Title:    "List of Products",
			Heading:  "Products",
			Products: products,
		}
		return c.Render("index", data)
	}
}
