package server

import (
	"encoding/json"
	"hippias-fiber/internal/models"
	_ "hippias-fiber/swagger"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/supabase-community/supabase-go"
)

type Server struct {
	*fiber.App
	sb *supabase.Client
}

func New() *Server {
	API_KEY := os.Getenv("API_KEY")
	API_URL := os.Getenv("API_URL")
	client, err := supabase.NewClient(API_URL, API_KEY, nil)
	if err != nil {
		panic(err)
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET",
	}))
	app.Get("/swagger/*", swagger.HandlerDefault)
	server := &Server{
		App: app,
		sb:  client,
	}

	server.setupRoutes()

	return server
}

func (s *Server) Sb() *supabase.Client {
	return s.sb
}

func (s *Server) setupRoutes() {
	s.App.Get("/book/id", s.getBook)
	s.App.Get("/list", s.listBooks)
	s.App.Get("/authors", s.listAuthors)
	s.App.Get("/authors/:id", s.getAuthor)
	s.App.Get("/books/author/:id", s.getBooksByAuthorID)
	s.App.Get("/courses", s.listCourses)
	s.App.Get("/courses/:id", s.getCourse)
	s.App.Get("/courses/details/:id", s.GetCourseWithDetails)
	s.App.Post("/courses", s.createCourse)
	s.App.Get("/facilitators", s.listFacilitators)
	s.App.Get("/facilitators/:id", s.getFacilitator)
	s.App.Post("/facilitators", s.createFacilitator)
	s.App.Delete("/facilitators/:id", s.deleteFacilitator)
}

// getBook godoc
// @Summary Get a book by ID
// @Description Retrieves a book by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} models.Book
// @Failure 500 {object} ErrorResponse
// @Router /book/{id} [get]
func (s *Server) getBook(c *fiber.Ctx) error {
	bookID := c.Params("id")

	data, _, err := s.sb.From("books").
		Select("*", "1", false).
		Eq("id", bookID).
		Single().
		Execute()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var book models.Book
	if err := json.Unmarshal(data, &book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	return c.JSON(book)
}

// getBooksByAuthorID godoc
// @Summary Get books by author ID
// @Description Retrieves books by the author's ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Author ID"
// @Success 200 {array} models.Book
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /books/author/{id} [get]
func (s *Server) getBooksByAuthorID(c *fiber.Ctx) error {
	authorID := c.Params("id")
	log.Printf("Author ID: %v", authorID)

	if authorID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: "Missing author ID"})
	}

	data, _, err := s.sb.From("books").
		Select("*", "exact", false).
		Eq("authorId", authorID).
		Execute()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var books []models.Book
	if err := json.Unmarshal(data, &books); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}
	log.Printf("Books: %v", books)

	return c.JSON(books)
}

// listAuthors godoc
// @Summary List authors
// @Description Retrieves a list of authors
// @Tags authors
// @Accept json
// @Produce json
// @Success 200 {array} models.Author
// @Failure 500 {object} ErrorResponse
// @Router /authors [get]
func (s *Server) listAuthors(c *fiber.Ctx) error {
	data, _, err := s.sb.From("authors").Select("*", "exact", false).Execute()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var authors []models.Author
	if err := json.Unmarshal(data, &authors); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	return c.JSON(authors)
}

// getAuthor godoc
// @Summary Get an author by ID
// @Description Retrieves an author by their ID
// @Tags authors
// @Accept json
// @Produce json
// @Param id path int true "Author ID"
// @Success 200 {object} models.Author
// @Failure 500 {object} ErrorResponse
// @Router /authors/{id} [get]
func (s *Server) getAuthor(c *fiber.Ctx) error {
	authorID := c.Params("id")
	log.Printf("Author ID: %v", authorID)
	data, _, err := s.sb.From("authors").
		Select("*", "exact", false).
		Eq("id", authorID).
		Single().
		Execute()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var author models.Author
	if err := json.Unmarshal(data, &author); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	return c.JSON(author)
}

// listBooks godoc
// @Summary List books
// @Description Retrieves a list of books
// @Tags books
// @Accept json
// @Produce json
// @Success 200 {array} models.Book
// @Failure 500 {object} ErrorResponse
// @Router /list [get]
func (s *Server) listBooks(c *fiber.Ctx) error {
	data, _, err := s.sb.From("books").Select("*", "exact", false).Execute()

	if err != nil {
		log.Printf("Error querying books: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var books []models.Book
	if err := json.Unmarshal(data, &books); err != nil {
		log.Printf("Error unmarshaling books: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Books: %+v", books)

	return c.JSON(books)
}

// listCourses godoc
// @Summary List courses
// @Description Retrieves a list of courses
// @Tags courses
// @Accept json
// @Produce json
// @Success 200 {array} models.Course
// @Failure 500 {object} ErrorResponse
// @Router /courses [get]
func (s *Server) listCourses(c *fiber.Ctx) error {
	data, _, err := s.sb.From("courses").Select("*", "exact", false).Execute()
	if err != nil {
		log.Printf("Error querying courses: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var courses []models.Course
	if err := json.Unmarshal(data, &courses); err != nil {
		log.Printf("Error unmarshaling courses: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Courses: %+v", courses)
	return c.JSON(courses)
}

// getCourse godoc
// @Summary Get a course by ID
// @Description Retrieves a course by its ID
// @Tags courses
// @Accept json
// @Produce json
// @Param id path int true "Course ID"
// @Success 200 {object} models.Course
// @Failure 500 {object} ErrorResponse
// @Router /courses/{id} [get]
func (s *Server) getCourse(c *fiber.Ctx) error {
	courseID := c.Params("id")

	data, _, err := s.Sb().From("courses").
		Select("*", "exact", false).
		Eq("id", courseID).
		Single().
		Execute()
	if err != nil {
		log.Printf("Error querying course: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var course models.Course
	if err := json.Unmarshal(data, &course); err != nil {
		log.Printf("Error unmarshaling course: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Course: %+v", course)
	return c.JSON(course)
}

// GetCourseWithDetails godoc
// @Summary Get course details with facilitator and books
// @Description Retrieves the course details along with its associated facilitator and an array of books included in the course
// @Tags courses
// @Accept json
// @Produce json
// @Param id path int true "Course ID"
// @Success 200 {object} GetCourseWithDetailsResponse
// @Failure 500 {object} ErrorResponse
// @Router /courses/details/{id} [get]
func (s *Server) GetCourseWithDetails(c *fiber.Ctx) error {
	courseID := c.Params("id")

	courseData, _, err := s.Sb().From("courses").
		Select("*", "exact", false).
		Eq("id", courseID).
		Single().
		Execute()
	if err != nil {
		log.Printf("Error querying course: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var course models.Course
	if err := json.Unmarshal(courseData, &course); err != nil {
		log.Printf("Error unmarshaling course: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}
	facIdStr := strconv.Itoa(course.FacilitatorID)
	facilitatorData, _, err := s.Sb().From("facilitators").
		Select("*", "exact", false).
		Eq("id", facIdStr).
		Single().
		Execute()
	if err != nil {
		log.Printf("Error querying facilitator: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var facilitator models.Facilitator
	if err := json.Unmarshal(facilitatorData, &facilitator); err != nil {
		log.Printf("Error unmarshaling facilitator: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	booksData, _, err := s.Sb().From("course_books").
		Select("book_id", "exact", false).
		Eq("course_id", courseID).
		Execute()
	if err != nil {
		log.Printf("Error querying course books: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var courseBooks []models.CourseBook
	if err := json.Unmarshal(booksData, &courseBooks); err != nil {
		log.Printf("Error unmarshaling course books: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var books []models.Book

	for _, courseBook := range courseBooks {
		var bId = strconv.Itoa(courseBook.BookID)
		bookData, _, err := s.Sb().From("books").
			Select("*", "exact", false).
			Eq("id", bId).
			Single().
			Execute()
		if err != nil {
			log.Printf("Error querying book: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
		}

		var book models.Book
		if err := json.Unmarshal(bookData, &book); err != nil {
			log.Printf("Error unmarshaling book: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
		}

		books = append(books, book)
	}

	response := GetCourseWithDetailsResponse{
		Course:      course,
		Facilitator: facilitator,
		Books:       books,
	}

	return c.JSON(response)
}

type GetCourseWithDetailsResponse struct {
	Course      models.Course      `json:"course"`
	Facilitator models.Facilitator `json:"facilitator"`
	Books       []models.Book      `json:"books"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

// createCourse godoc
// @Summary Create a course
// @Description Creates a new course
// @Tags courses
// @Accept json
// @Produce json
// @Param course body models.Course true "Course object"
// @Success 200 {object} models.Course
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /courses [post]
func (s *Server) createCourse(c *fiber.Ctx) error {
	var course models.Course
	if err := c.BodyParser(&course); err != nil {
		log.Printf("Error parsing course: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: err.Error()})
	}

	data, err := json.Marshal(course)
	if err != nil {
		log.Printf("Error marshaling course: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	data, _, err = s.sb.From("courses").Insert(string(data), false, "Error", "Success", "1").Execute()
	if err != nil {
		log.Printf("Error inserting course: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Created course: %+v", course)
	return c.JSON(data)
}

// listFacilitators godoc
// @Summary List facilitators
// @Description Retrieves a list of facilitators
// @Tags facilitators
// @Accept json
// @Produce json
// @Success 200 {array} models.Facilitator
// @Failure 500 {object} ErrorResponse
// @Router /facilitators [get]
func (s *Server) listFacilitators(c *fiber.Ctx) error {
	data, _, err := s.sb.From("facilitators").Select("*", "exact", true).Execute()
	if err != nil {
		log.Printf("Error querying facilitators: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var facilitators []models.Facilitator
	if err := json.Unmarshal(data, &facilitators); err != nil {
		log.Printf("Error unmarshaling facilitators: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Facilitators: %+v", facilitators)
	return c.JSON(facilitators)
}

// getFacilitator godoc
// @Summary Get a facilitator by ID
// @Description Retrieves a facilitator by their ID
// @Tags facilitators
// @Accept json
// @Produce json
// @Param id path int true "Facilitator ID"
// @Success 200 {object} models.Facilitator
// @Failure 500 {object} ErrorResponse
// @Router /facilitators/{id} [get]
func (s *Server) getFacilitator(c *fiber.Ctx) error {
	facilitatorID := c.Params("id")

	data, _, err := s.sb.From("facilitators").Select("*", "exact", true).Eq("id", facilitatorID).Execute()
	if err != nil {
		log.Printf("Error querying facilitator: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var facilitator models.Facilitator
	if err := json.Unmarshal(data, &facilitator); err != nil {
		log.Printf("Error unmarshaling facilitator: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Facilitator: %+v", facilitator)
	return c.JSON(facilitator)
}

// createFacilitator godoc
// @Summary Create a facilitator
// @Description Creates a new facilitator
// @Tags facilitators
// @Accept json
// @Produce json
// @Param facilitator body models.Facilitator true "Facilitator object"
// @Success 200 {object} models.Facilitator
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /facilitators [post]
func (s *Server) createFacilitator(c *fiber.Ctx) error {
	var facilitator models.Facilitator
	if err := c.BodyParser(&facilitator); err != nil {
		log.Printf("Error parsing facilitator: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: err.Error()})
	}

	_, _, err := s.sb.From("facilitators").Insert(facilitator, false, "", "*", "").Execute()
	if err != nil {
		log.Printf("Error inserting facilitator: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Created facilitator: %+v", facilitator)
	return c.JSON(facilitator)
}

// deleteFacilitator godoc
// @Summary Delete a facilitator by ID
// @Description Deletes a facilitator by their ID
// @Tags facilitators
// @Accept json
// @Produce json
// @Param id path int true "Facilitator ID"
// @Success 204
// @Failure 500 {object} ErrorResponse
// @Router /facilitators/{id} [delete]
func (s *Server) deleteFacilitator(c *fiber.Ctx) error {
	facilitatorID := c.Params("id")

	_, _, err := s.sb.From("facilitators").Delete("Success", "true").Eq("id", facilitatorID).Execute()
	if err != nil {
		log.Printf("Error deleting facilitator: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Deleted facilitator with ID: %s", facilitatorID)
	return c.SendStatus(fiber.StatusNoContent)
}
