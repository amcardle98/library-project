package main

import (
	"challenge/web-service-gin/models"
	"challenge/web-service-gin/storage"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

type Book struct {
	Title string `json: "title"`
}

func (r *Repository) GetBookByTitle(c *fiber.Ctx) error {
	title := c.Params("title")
	bookModel := &models.Book{}
	if title == "" {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Title is required",
		})
		return nil
	}

	fmt.Println(title)

	err := r.DB.Where("lower(title) LIKE '%' || ? || '%'", strings.ToLower(title)).First(&bookModel).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Book not found",
		})
		return err
	}
	c.Status(http.StatusOK).JSON(&fiber.Map{
		"book": bookModel,
	})
	return nil
}

func (r *Repository) GetBookById(c *fiber.Ctx) error {
	id := c.Params("id")
	bookModel := &models.Book{}
	if id == "" {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "ID is required",
		})
		return nil
	}

	fmt.Println(id)

	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Book not found",
		})
		return err
	}
	c.Status(http.StatusOK).JSON(&fiber.Map{
		"data": bookModel,
	})
	return nil
}

func (r *Repository) GetBooks(c *fiber.Ctx) error {
	bookModels := &[]models.Book{}

	err := r.DB.Find(bookModels).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Books not found",
		})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Books found",
		"data":    bookModels,
	})
	return nil
}

func (r *Repository) CreateBook(c *fiber.Ctx) error {
	book := Book{}

	err := c.BodyParser(&book)
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Book not created",
		})
		return err
	}

	err = r.DB.Create(&book).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Book not created",
		})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Book created"})
	return nil

	//return nil
}

func (r *Repository) UpdateBook(c *fiber.Ctx) error {
	bookModel := new(Book)
	if err := c.BodyParser(bookModel); err != nil {
		return err
	}
	err := r.DB.Save(bookModel).Error
	return err

	//return nil
}

func (r *Repository) DeleteBook(c *fiber.Ctx) error {
	bookModel := models.Book{}
	id := c.Params("id")
	if id == "" {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "ID is required",
		})
		return nil
	}

	err := r.DB.Delete(bookModel, id).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Book not deleted",
		})
		return err
	}
	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Book deleted",
	})
	return nil
}

func (r *Repository) InitDB(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/getByTitle/:title", r.GetBookByTitle)
	api.Get("/books", r.GetBooks)
	api.Get("/getById/:id", r.GetBookById)
	api.Post("/create", r.CreateBook)
	api.Delete("/delete_book/:id", r.DeleteBook)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal(err)
	}

	err = models.Migration(db)
	if err != nil {
		log.Fatal(err)
	}

	r := Repository{
		DB: db,
	}

	// Connect to the database
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
	}))

	r.InitDB(app)
	app.Listen(":4000")
}
