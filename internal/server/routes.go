package server

import (
	"encoding/json"
	"log"

	"hippias-fiber/internal/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *Server) RegisterRoutes() *fiber.App {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, DELETE",
	}))

	// Book routes
	app.Get("/books/:id", handlers.GetBook(s))
	app.Get("/books", handlers.ListBooks(s))
	app.Get("/books/author/:id", handlers.GetBooksByAuthorID(s))

	// Author routes
	app.Get("/authors", handlers.ListAuthors(s))
	app.Get("/authors/:id", handlers.GetAuthor(s))
	app.Post("/authors", handlers.CreateAuthor(s))
	app.Delete("/authors/:id", handlers.DeleteAuthor(s))

	// Course routes
	app.Get("/courses", handlers.ListCourses(s))
	app.Get("/courses/:id", handlers.GetCourse(s))
	app.Post("/courses", handlers.CreateCourse(s))

	// Facilitator routes
	app.Get("/facilitators", handlers.ListFacilitators(s))
	app.Get("/facilitators/:id", handlers.GetFacilitator(s))
	app.Post("/facilitators", handlers.CreateFacilitator(s))
	app.Delete("/facilitators/:id", handlers.DeleteFacilitator(s))

	// CourseParticipant routes
	app.Get("/course-participants", handlers.ListCourseParticipants(s))
	app.Get("/course-participants/:courseId/:userId", handlers.GetCourseParticipant(s))
	app.Post("/course-participants", handlers.CreateCourseParticipant(s))
	app.Delete("/course-participants/:courseId/:userId", handlers.DeleteCourseParticipant(s))

	return app
}

/*
*
__

	     ____              _
		| __ )  ___   ___ | | __
		|  _ \ / _ \ / _ \| |/ /
		| |_) | (_) | (_) |   <
		|____/ \___/ \___/|_|\_\
		*
*/
/* func (s *Server) listBooks(c *fiber.Ctx) error {
	data, _, err := s.Sb.From("books").Select("*", "exact", false).Execute()

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
*/
/*
*	 ____
	/ ___|___  _   _ _ __ ___  ___  ___
    | |   / _ \| | | | '__/ __|/ _ \/ __|
	| |__| (_) | |_| | |  \__ \  __/\__ \
	 \____\___/ \__,_|_|  |___/\___||___/
	 *
*/

func (s *Server) listCourses(c *fiber.Ctx) error {
	data, _, err := s.Sb.From("courses").Select("*", "exact", false).Execute()
	if err != nil {
		log.Printf("Error querying courses: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var courses []Course
	if err := json.Unmarshal(data, &courses); err != nil {
		log.Printf("Error unmarshaling courses: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	log.Printf("Courses: %+v", courses)
	return c.JSON(courses)
}

func (s *Server) getCourse(c *fiber.Ctx) error {
	courseID := c.Params("id")

	data, _, err := s.Sb.From("courses").Select("*", "exact", false).Eq("id", courseID).Execute()
	if err != nil {
		log.Printf("Error querying course: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var course Course
	if err := json.Unmarshal(data, &course); err != nil {
		log.Printf("Error unmarshaling course: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	log.Printf("Course: %+v", course)
	return c.JSON(course)
}

func (s *Server) createCourse(c *fiber.Ctx) error {
	var course Course
	if err := c.BodyParser(&course); err != nil {
		log.Printf("Error parsing course: %v", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	data, err := json.Marshal(course)
	if err != nil {
		log.Printf("Error marshaling course: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	data, _, err = s.Sb.From("courses").Insert(string(data), false, "Error", "Success", "1").Execute()
	if err != nil {
		log.Printf("Error inserting course: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	log.Printf("Created course: %+v", course)
	return c.JSON(data)
}

/**
_____           _ _ _ _        _
 |  ___|_ _  ___(_) (_) |_ __ _| |_ ___  _ __
 | |_ / _` |/ __| | | | __/ _` | __/ _ \| '__|
 |  _| (_| | (__| | | | || (_| | || (_) | |
 |_|  \__,_|\___|_|_|_|\__\__,_|\__\___/|_|
**/

func (s *Server) listFacilitators(c *fiber.Ctx) error {
	data, _, err := s.Sb.From("facilitators").Select("*", "exact", true).Execute()
	if err != nil {
		log.Printf("Error querying facilitators: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var facilitators []Facilitator
	if err := json.Unmarshal(data, &facilitators); err != nil {
		log.Printf("Error unmarshaling facilitators: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	log.Printf("Facilitators: %+v", facilitators)
	return c.JSON(facilitators)
}

func (s *Server) getFacilitator(c *fiber.Ctx) error {
	facilitatorID := c.Params("id")

	data, _, err := s.Sb.From("facilitators").Select("*", "exact", true).Eq("id", facilitatorID).Execute()
	if err != nil {
		log.Printf("Error querying facilitator: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var facilitator Facilitator
	if err := json.Unmarshal(data, &facilitator); err != nil {
		log.Printf("Error unmarshaling facilitator: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	log.Printf("Facilitator: %+v", facilitator)
	return c.JSON(facilitator)
}

func (s *Server) createFacilitator(c *fiber.Ctx) error {
	var facilitator Facilitator
	if err := c.BodyParser(&facilitator); err != nil {
		log.Printf("Error parsing facilitator: %v", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	data, _, err := s.Sb.From("facilitators").Insert(facilitator, false, "", "*", "").Execute()
	if err != nil {
		log.Printf("Error inserting facilitator: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	log.Printf("Created facilitator: %+v", facilitator, data)
	return c.JSON(facilitator)
}

func (s *Server) deleteFacilitator(c *fiber.Ctx) error {
	facilitatorID := c.Params("id")

	data, _, err := s.Sb.From("facilitators").Delete("Success", "true").Eq("id", facilitatorID).Execute()
	if err != nil {
		log.Printf("Error deleting facilitator: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	log.Printf("Deleted facilitator with ID: %s", facilitatorID, data)
	return c.SendStatus(fiber.StatusNoContent)
}

/**
____                          ____              _
  / ___|___  _   _ _ __ ___  ___|  _ \  ___   ___ | | __
 | |   / _ \| | | | '__/ __|/ _ \ | | |/ _ \ / _ \| |/ /
 | |__| (_) | |_| | |  \__ \  __/ |_| | (_) | (_) |   <
  \____\___/ \__,_|_|  |___/\___|____/ \___/ \___/|_|\_\
  **/
