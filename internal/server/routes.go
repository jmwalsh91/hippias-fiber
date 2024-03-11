package server

import (
	"encoding/json"
	"log"
	"time"

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
	app.Get("/courses", s.listCourses)
	app.Get("/courses/:id", s.getCourse)
	app.Post("/courses", s.createCourse)
	app.Get("/facilitators", s.listFacilitators)
	app.Get("/facilitators/:id", s.getFacilitator)
	app.Post("/facilitators", s.createFacilitator)
	app.Delete("/facilitators/:id", s.deleteFacilitator)
	return app
}

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Facilitator struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Course struct {
	ID            int       `json:"id"`
	FacilitatorID int       `json:"facilitatorId"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Book struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	AuthorID    int       `json:"authorId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Author struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Nationality string    `json:"nationality"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CourseBook struct {
	ID        int       `json:"id"`
	CourseID  int       `json:"courseId"`
	BookID    int       `json:"bookId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CourseParticipant struct {
	ID        int       `json:"id"`
	CourseID  int       `json:"courseId"`
	UserID    int       `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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

/*
*

		_         _   _
	   / \  _   _| |_| |__   ___  _ __
	  / _ \| | | | __| '_ \ / _ \| '__|
	 / ___ \ |_| | |_| | | | (_) | |
	/_/   \_\__,_|\__|_| |_|\___/|_|

*
*/
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

/*
*	 ____
	/ ___|___  _   _ _ __ ___  ___  ___
    | |   / _ \| | | | '__/ __|/ _ \/ __|
	| |__| (_) | |_| | |  \__ \  __/\__ \
	 \____\___/ \__,_|_|  |___/\___||___/
	 *
*/

func (s *Server) listCourses(c *fiber.Ctx) error {
	data, _, err := s.sb.From("courses").Select("*", "exact", false).Execute()
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

	data, _, err := s.sb.From("courses").Select("*", "exact", false).Eq("id", courseID).Execute()
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

	data, _, err = s.sb.From("courses").Insert(string(data), false, "Error", "Success", "1").Execute()
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
	data, _, err := s.sb.From("facilitators").Select("*", "exact", true).Execute()
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

	data, _, err := s.sb.From("facilitators").Select("*", "exact", true).Eq("id", facilitatorID).Execute()
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

	data, _, err := s.sb.From("facilitators").Insert(facilitator, false, "", "*", "").Execute()
	if err != nil {
		log.Printf("Error inserting facilitator: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	log.Printf("Created facilitator: %+v", facilitator, data)
	return c.JSON(facilitator)
}

func (s *Server) deleteFacilitator(c *fiber.Ctx) error {
	facilitatorID := c.Params("id")

	data, _, err := s.sb.From("facilitators").Delete("Success", "true").Eq("id", facilitatorID).Execute()
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
