package handlers

import (
	"encoding/json"
	"log"

	"hippias-fiber/internal/models"
	"hippias-fiber/internal/server"

	"github.com/gofiber/fiber/v2"
)

func ListCourseBooks(s *server.Server) fiber.Handler {
	return func(c *fiber.Ctx) error {
		data, _, err := s.Sb.From("course_books").
			Select("*", "exact", true).
			Execute()
		if err != nil {
			log.Printf("Error querying course books: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		var courseBooks []models.CourseBook
		if err := json.Unmarshal(data, &courseBooks); err != nil {
			log.Printf("Error unmarshaling course books: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Printf("Course Books: %v", courseBooks)
		return c.JSON(courseBooks)
	}
}

func GetCourseBook(s *server.Server) fiber.Handler {
	return func(c *fiber.Ctx) error {
		courseID := c.Params("courseId")
		bookID := c.Params("bookId")

		data, _, err := s.Sb.From("course_books").
			Select("*", "exact", true).
			Eq("course_id", courseID).
			Eq("book_id", bookID).
			Execute()
		if err != nil {
			log.Printf("Error querying course book: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		var courseBook models.CourseBook
		if err := json.Unmarshal(data, &courseBook); err != nil {
			log.Printf("Error unmarshaling course book: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Printf("Course Book: %v", courseBook)
		return c.JSON(courseBook)
	}
}

func CreateCourseBook(s *server.Server) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var courseBook models.CourseBook
		if err := c.BodyParser(&courseBook); err != nil {
			log.Printf("Error parsing course book: %v", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		_, _, err := s.Sb.From("course_books").
			Insert(courseBook, false, "", "*", "").
			Execute()
		if err != nil {
			log.Printf("Error inserting course book: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Printf("Created course book: %v", courseBook)
		return c.JSON(courseBook)
	}
}

func DeleteCourseBook(s *server.Server) fiber.Handler {
	return func(c *fiber.Ctx) error {
		courseID := c.Params("courseId")
		bookID := c.Params("bookId")

		_, _, err := s.Sb.From("course_books").
			Delete("Success", "true").
			Eq("course_id", courseID).
			Eq("book_id", bookID).
			Execute()
		if err != nil {
			log.Printf("Error deleting course book: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Printf("Deleted course book with course ID: %s and book ID: %s", courseID, bookID)
		return c.SendStatus(fiber.StatusNoContent)
	}
}
