package handlers

import (
	"encoding/json"
	"log"

	"hippias-fiber/internal/models"
	"hippias-fiber/internal/server"

	"github.com/gofiber/fiber/v2"
)

func ListCourseParticipants(s *server.Server) fiber.Handler {
	return func(c *fiber.Ctx) error {
		data, _, err := s.Sb.From("course_participants").
			Select("*", "exact", true).
			Execute()
		if err != nil {
			log.Printf("Error querying course participants: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		var courseParticipants []models.CourseParticipant
		if err := json.Unmarshal(data, &courseParticipants); err != nil {
			log.Printf("Error unmarshaling course participants: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Printf("Course Participants: %v", courseParticipants)
		return c.JSON(courseParticipants)
	}
}

func GetCourseParticipant(s *server.Server) fiber.Handler {
	return func(c *fiber.Ctx) error {
		courseID := c.Params("courseId")
		userID := c.Params("userId")

		data, _, err := s.Sb.From("course_participants").
			Select("*", "exact", true).
			Eq("course_id", courseID).
			Eq("user_id", userID).
			Execute()
		if err != nil {
			log.Printf("Error querying course participant: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		var courseParticipant models.CourseParticipant
		if err := json.Unmarshal(data, &courseParticipant); err != nil {
			log.Printf("Error unmarshaling course participant: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Printf("Course Participant: %v", courseParticipant)
		return c.JSON(courseParticipant)
	}
}

func CreateCourseParticipant(s *server.Server) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var courseParticipant models.CourseParticipant
		if err := c.BodyParser(&courseParticipant); err != nil {
			log.Printf("Error parsing course participant: %v", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		_, _, err := s.Sb.From("course_participants").
			Insert(courseParticipant, false, "", "*", "").
			Execute()
		if err != nil {
			log.Printf("Error inserting course participant: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Printf("Created course participant: %v", courseParticipant)
		return c.JSON(courseParticipant)
	}
}

func DeleteCourseParticipant(s *server.Server) fiber.Handler {
	return func(c *fiber.Ctx) error {
		courseID := c.Params("courseId")
		userID := c.Params("userId")

		_, _, err := s.Sb.From("course_participants").
			Delete("Success", "true").
			Eq("course_id", courseID).
			Eq("user_id", userID).
			Execute()
		if err != nil {
			log.Printf("Error deleting course participant: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Printf("Deleted course participant with course ID: %s and user ID: %s", courseID, userID)
		return c.SendStatus(fiber.StatusNoContent)
	}
}
