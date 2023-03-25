package handlers

import (
	"template/database"
	"template/models"

	"github.com/gofiber/fiber/v2"
)

type PageData struct {
	Title    string
	Heading  string
	Products []models.Product
}

type SinglePageData struct {
	Title   string
	Heading string
	Product models.Product
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

	return c.Render("product", data)
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
		database.DB.Db.Where("name LIKE ?", "%"+searchQuery+"%").Limit(100).Find(&products)
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
