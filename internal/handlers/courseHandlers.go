package handlers

import (
	"encoding/json"
	"hippias-fiber/internal/models"
	"hippias-fiber/internal/server"
	"log"

	"github.com/gofiber/fiber/v2"
)

/*
		 ____
		/ ___|___  _   _ _ __ ___  ___  ___
	    | |   / _ \| | | | '__/ __|/ _ \/ __|
		| |__| (_) | |_| | |  \__ \  __/\__ \
		 \____\___/ \__,_|_|  |___/\___||___/
		 *
*/
func ListCourses(s *server.Server) fiber.Handler {
	return func(c *fiber.Ctx) error {
		data, _, err := s.Sb.From("courses").
			Select("*", "exact", false).
			Execute()
		if err != nil {
			log.Printf("Error querying courses: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		var courses []models.Course
		if err := json.Unmarshal(data, &courses); err != nil {
			log.Printf("Error unmarshaling courses: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Printf("Courses: %+v", courses)
		return c.JSON(courses)
	}
}

func GetCourse(s *server.Server) fiber.Handler {
	return func(c *fiber.Ctx) error {
		courseID := c.Params("id")

		data, _, err := s.Sb.From("courses").
			Select("*", "exact", false).
			Eq("id", courseID).
			Execute()
		if err != nil {
			log.Printf("Error querying course: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		var course models.Course
		if err := json.Unmarshal(data, &course); err != nil {
			log.Printf("Error unmarshaling course: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Printf("Course: %+v", course)
		return c.JSON(course)
	}
}

func CreateCourse(s *server.Server) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var course models.Course
		if err := c.BodyParser(&course); err != nil {
			log.Printf("Error parsing course: %v", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		data, err := json.Marshal(course)
		if err != nil {
			log.Printf("Error marshaling course: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		data, _, err = s.Sb.From("courses").
			Insert(string(data), false, "Error", "Success", "1").
			Execute()
		if err != nil {
			log.Printf("Error inserting course: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Printf("Created course: %+v", course)
		return c.JSON(data)
	}
}
