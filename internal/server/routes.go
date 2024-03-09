package server

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *Server) RegisterRoutes() *fiber.App {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET",
	}))

	app.Get("/book/id", s.getBook)
	app.Get("/list", s.listBooks)
	app.Get("/authors", s.listAuthors)
	app.Get("/authors/:id", s.getAuthor)
	app.Get("/books/author/:id", s.getBooksByAuthorID)

	return app
}

type Book = struct {
	Id          int      `json:"id"`
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	Description string   `json:"description"`
	AuthorId    int      `json:"authorId"`
	Tags        []string `json:"tags"`
}

type Author struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Nationality string `json:"nationality"`
	Description string `json:"description"`
}

func (s *Server) getBook(c *fiber.Ctx) error {
	bookID := c.Params("id")

	data, _, err := s.sb.From("books").
		Select("*", "1", false).
		Eq("id", bookID).
		Single().
		Execute()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	var book Book
	if err := json.Unmarshal(data, &book); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(book)
}

func (s *Server) getBooksByAuthorID(c *fiber.Ctx) error {
	authorID := c.Params("id")
	log.Printf("Author ID: %v", authorID)

	if authorID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing author ID")
	}

	data, _, err := s.sb.From("books").
		Select("*", "exact", false).
		Eq("authorId", authorID).
		Execute()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	var books []Book
	if err := json.Unmarshal(data, &books); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	log.Printf("Books: %v", books)

	return c.JSON(books)
}

func (s *Server) listAuthors(c *fiber.Ctx) error {
	data, _, err := s.sb.From("authors").Select("*", "exact", false).Execute()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	var authors []Author
	if err := json.Unmarshal(data, &authors); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(authors)
}

func (s *Server) getAuthor(c *fiber.Ctx) error {
	authorID := c.Params("id")
	log.Printf("Author ID: %v", authorID)
	data, _, err := s.sb.From("authors").
		Select("*", "exact", false).
		Eq("id", authorID).
		Single().
		Execute()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	var author Author
	if err := json.Unmarshal(data, &author); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(author)
}

func (s *Server) listBooks(c *fiber.Ctx) error {
	data, _, err := s.sb.From("books").Select("*", "exact", false).Execute()

	if err != nil {
		log.Printf("Error querying books: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var books []Book
	if err := json.Unmarshal(data, &books); err != nil {
		log.Printf("Error unmarshaling books: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	log.Printf("Books: %+v", books)

	return c.JSON(books)
}
