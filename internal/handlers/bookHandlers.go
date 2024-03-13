package handlers

import (
	"encoding/json"
	"hippias-fiber/internal/models"
	"hippias-fiber/internal/server"
	"log"

	"github.com/gofiber/fiber/v2"
)

func GetBook(s *server.Server) fiber.Handler {
	return func(c *fiber.Ctx) error {
		bookID := c.Params("id")
		data, _, err := s.Sb.From("books").
			Select("*", "1", false).
			Eq("id", bookID).
			Single().
			Execute()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		var book models.Book
		if err := json.Unmarshal(data, &book); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.JSON(book)
	}
}

func GetBooksByAuthorID(s *server.Server) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorID := c.Params("id")
		log.Printf("Author ID: %v", authorID)

		if authorID == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Missing author ID")
		}

		data, _, err := s.Sb.From("books").
			Select("*", "exact", false).
			Eq("authorId", authorID).
			Execute()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		var books []models.Book
		if err := json.Unmarshal(data, &books); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		log.Printf("Books: %v", books)

		return c.JSON(books)

	}
}
