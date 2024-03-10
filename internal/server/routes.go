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

/*
*

	 ____
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

	_, err = s.sb.From("courses").Insert(string(data)).Execute()
	if err != nil {
		log.Printf("Error inserting course: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	log.Printf("Created course: %+v", course)
	return c.JSON(course)
}
