package handlers

import (
	"encoding/json"
	"hippias-fiber/internal/models"
	"hippias-fiber/internal/server"
	"log"

	"github.com/gofiber/fiber/v2"
)

func ListAuthors(s *server.Server) fiber.Handler {
	return func(c *fiber.Ctx) error {
		data, _, err := s.Sb.From("authors").
			Select("*", "exact", true).
			Execute()
		if err != nil {
			log.Printf("Error querying authors: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		var authors []models.Author
		if err := json.Unmarshal(data, &authors); err != nil {
			log.Printf("Error unmarshaling authors: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Printf("Authors: %v", authors)
		return c.JSON(authors)
	}
}

func GetAuthor(s *server.Server) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorID := c.Params("id")

		data, _, err := s.Sb.From("authors").
			Select("*", "exact", true).
			Eq("id", authorID).
			Execute()
		if err != nil {
			log.Printf("Error querying author: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		var author models.Author
		if err := json.Unmarshal(data, &author); err != nil {
			log.Printf("Error unmarshaling author: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Printf("Author: %v", author)
		return c.JSON(author)
	}
}

func CreateAuthor(s *server.Server) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var author models.Author
		if err := c.BodyParser(&author); err != nil {
			log.Printf("Error parsing author: %v", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		_, _, err := s.Sb.From("authors").
			Insert(author, false, "", "*", "").
			Execute()
		if err != nil {
			log.Printf("Error inserting author: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Printf("Created author: %v", author)
		return c.JSON(author)
	}
}

func DeleteAuthor(s *server.Server) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorID := c.Params("id")

		_, _, err := s.Sb.From("authors").
			Delete("Success", "true").
			Eq("id", authorID).
			Execute()
		if err != nil {
			log.Printf("Error deleting author: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Printf("Deleted author with ID: %s", authorID)
		return c.SendStatus(fiber.StatusNoContent)
	}
}
