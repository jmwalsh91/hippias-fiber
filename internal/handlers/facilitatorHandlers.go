package handlers

import (
	"encoding/json"
	"log"

	"hippias-fiber/internal/models"
	"hippias-fiber/internal/server"

	"github.com/gofiber/fiber/v2"
)

func ListFacilitators(s *server.Server) fiber.Handler {
	return func(c *fiber.Ctx) error {
		data, _, err := s.Sb.From("facilitators").
			Select("*", "exact", true).
			Execute()
		if err != nil {
			log.Printf("Error querying facilitators: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		var facilitators []models.Facilitator
		if err := json.Unmarshal(data, &facilitators); err != nil {
			log.Printf("Error unmarshaling facilitators: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Printf("Facilitators: %v", facilitators)
		return c.JSON(facilitators)
	}
}

func GetFacilitator(s *server.Server) fiber.Handler {
	return func(c *fiber.Ctx) error {
		facilitatorID := c.Params("id")

		data, _, err := s.Sb.From("facilitators").
			Select("*", "exact", true).
			Eq("id", facilitatorID).
			Execute()
		if err != nil {
			log.Printf("Error querying facilitator: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		var facilitator models.Facilitator
		if err := json.Unmarshal(data, &facilitator); err != nil {
			log.Printf("Error unmarshaling facilitator: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Printf("Facilitator: %v", facilitator)
		return c.JSON(facilitator)
	}
}

func CreateFacilitator(s *server.Server) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var facilitator models.Facilitator
		if err := c.BodyParser(&facilitator); err != nil {
			log.Printf("Error parsing facilitator: %v", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		_, _, err := s.Sb.From("facilitators").
			Insert(facilitator, false, "", "*", "").
			Execute()
		if err != nil {
			log.Printf("Error inserting facilitator: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Printf("Created facilitator: %v", facilitator)
		return c.JSON(facilitator)
	}
}

func DeleteFacilitator(s *server.Server) fiber.Handler {
	return func(c *fiber.Ctx) error {
		facilitatorID := c.Params("id")

		_, _, err := s.Sb.From("facilitators").
			Delete("Success", "true").
			Eq("id", facilitatorID).
			Execute()
		if err != nil {
			log.Printf("Error deleting facilitator: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Printf("Deleted facilitator with ID: %s", facilitatorID)
		return c.SendStatus(fiber.StatusNoContent)
	}
}
