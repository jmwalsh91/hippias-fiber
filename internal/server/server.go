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
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/swagger"
	"github.com/mitchellh/mapstructure"
	supa "github.com/nedpals/supabase-go"
)

type Server struct {
	*fiber.App
	sb *supa.Client
}

func getDecoder() *mapstructure.Decoder {
	config := &mapstructure.DecoderConfig{
		TagName: "json",
		Squash:  false,
		Result:  &map[string]interface{}{},
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		log.Fatalf("Error creating decoder: %v", err)
	}
	return decoder
}

func New() *Server {
	API_KEY := os.Getenv("API_KEY")
	API_URL := os.Getenv("API_URL")
	client := supa.CreateClient(API_URL, API_KEY)

	app := fiber.New()
	store := session.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("session", store)
		return c.Next()
	})
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

func (s *Server) Sb() *supa.Client {
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
	s.App.Post("/login", s.login)
	s.App.Post("/register", s.register)
	s.App.Post("/logout", s.logout)
	s.App.Get("/discussions", s.listDiscussions)
	s.App.Get("/discussions/:id", s.getDiscussion)
	s.App.Post("/discussions", s.createDiscussion)
	s.App.Put("/discussions/:id", s.updateDiscussion)
	s.App.Delete("/discussions/:id", s.deleteDiscussion)
	s.App.Post("/reading-ratings", s.createReadingRating)
	s.App.Get("/reading-ratings/:id", s.getReadingRating)
	s.App.Get("/readings/:id/ratings", s.listReadingRatings)
	s.App.Put("/reading-ratings/:id", s.updateReadingRating)
	s.App.Delete("/reading-ratings/:id", s.deleteReadingRating)
	s.App.Get("/readings", s.listReadings)
	s.App.Get("/readings/:id", s.getReading)
	s.App.Post("/readings", s.createReading)
	s.App.Put("/readings/:id", s.updateReading)
	s.App.Delete("/readings/:id", s.deleteReading)
	s.App.Post("/discussion-attendance", s.createDiscussionAttendance)
	s.App.Get("/discussions/:id/attendance", s.listDiscussionAttendance)
	s.App.Get("/courses/:id/management", s.getCourseManagementDetails)
	s.App.Get("/discussions/:id/management", s.GetDiscussionMgmtDetails)
}

func (s *Server) login(c *fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: err.Error()})
	}

	user, err := s.sb.Auth.SignIn(c.Context(), supa.UserCredentials{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}
	log.Printf("User: %+v", user)
	return c.JSON(map[string]string{"message": "Login successful"})
}

func (s *Server) logout(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	err := s.sb.Auth.SignOut(c.Context(), token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}
	return c.JSON(map[string]string{"message": "Logout successful"})
}

func (s *Server) register(c *fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: err.Error()})
	}
	user, err := s.sb.Auth.SignUp(c.Context(), supa.UserCredentials{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}
	log.Printf("User: %+v", user)
	return c.JSON(map[string]string{"message": "Registration successful"})
}
func (s *Server) getBook(c *fiber.Ctx) error {
	bookID := c.Params("id")

	var result map[string]interface{}
	err := s.sb.DB.From("books").
		Select("*").
		Single().
		Eq("id", bookID).
		Execute(&result)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var book models.Book
	if err := mapstructure.Decode(result, &book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	return c.JSON(book)
}

func (s *Server) getBooksByAuthorID(c *fiber.Ctx) error {
	authorID := c.Params("id")
	log.Printf("Author ID: %v", authorID)

	if authorID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: "Missing author ID"})
	}

	var result map[string]interface{}
	err := s.sb.DB.From("books").
		Select("*").
		Eq("authorId", authorID).
		Execute(&result)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Result: %v", result)
	var books []models.Book
	if err := mapstructure.Decode(result, &books); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}
	log.Printf("Books: %v", books)

	return c.JSON(books)
}

func (s *Server) listAuthors(c *fiber.Ctx) error {
	var results []map[string]interface{}
	err := s.sb.DB.From("authors").Select("*").Execute(&results)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var authors []models.Author
	if err := mapstructure.Decode(results, &authors); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	return c.JSON(authors)
}

func (s *Server) getAuthor(c *fiber.Ctx) error {
	authorID := c.Params("id")
	log.Printf("Author ID: %v", authorID)
	var result map[string]interface{}
	err := s.sb.DB.From("authors").
		Select("*").
		Single().
		Eq("id", authorID).
		Execute(&result)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var author models.Author
	if err := mapstructure.Decode(result, &author); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	return c.JSON(author)
}

func (s *Server) listBooks(c *fiber.Ctx) error {
	var result []map[string]interface{}
	err := s.sb.DB.From("books").Select("*").Execute(&result)
	if err != nil {
		log.Printf("Error querying books: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var books []models.Book
	if err := mapstructure.Decode(result, &books); err != nil {
		log.Printf("Error unmarshaling books: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Books: %+v", books)

	return c.JSON(books)
}

func (s *Server) listCourses(c *fiber.Ctx) error {
	var jsonResult json.RawMessage
	err := s.sb.DB.From("courses").Select("*").Execute(&jsonResult)
	if err != nil {
		log.Printf("Error querying courses: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var courses []models.Course
	err = json.Unmarshal(jsonResult, &courses)
	if err != nil {
		log.Printf("Error unmarshaling courses: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Courses: %+v", courses)
	return c.JSON(courses)
}
func (s *Server) getCourse(c *fiber.Ctx) error {
	courseID := c.Params("id")

	var jsonResult json.RawMessage
	err := s.sb.DB.From("courses").
		Select("*").
		Single().
		Eq("id", courseID).
		Execute(&jsonResult)
	if err != nil {
		log.Printf("Error querying course: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var course models.Course
	err = json.Unmarshal(jsonResult, &course)
	if err != nil {
		log.Printf("Error unmarshaling course: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Course: %+v", course)
	return c.JSON(course)
}

func (s *Server) GetCourseWithDetails(c *fiber.Ctx) error {
	courseID := c.Params("id")
	log.Printf("GetCourseWithDetails: Processing request for course ID: %s", courseID)

	var jsonResult json.RawMessage
	err := s.sb.DB.From("courses").
		Select("*").
		Single().
		Eq("id", courseID).
		Execute(&jsonResult)
	if err != nil {
		log.Printf("GetCourseWithDetails: Error querying course: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var course models.Course
	err = json.Unmarshal(jsonResult, &course)
	if err != nil {
		log.Printf("GetCourseWithDetails: Error unmarshaling course data: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}
	log.Printf("GetCourseWithDetails: Fetched course: %+v", course)

	if course.FacilitatorID == 0 {
		log.Printf("GetCourseWithDetails: Invalid Facilitator ID for course ID: %s", courseID)
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Message: "Facilitator not found"})
	}

	facId := strconv.Itoa(course.FacilitatorID)
	var jsonFacilitatorResult json.RawMessage
	err2 := s.sb.DB.From("facilitators").
		Select("*").
		Single().
		Eq("id", facId).
		Execute(&jsonFacilitatorResult)
	if err2 != nil {
		log.Printf("GetCourseWithDetails: Error querying facilitator: %v", err)
		log.Printf("GetCourseWithDetails: Facilitator ID: %d", course.FacilitatorID)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var facilitator models.Facilitator
	err = json.Unmarshal(jsonFacilitatorResult, &facilitator)
	if err != nil {
		log.Printf("GetCourseWithDetails: Error unmarshaling facilitator data: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}
	log.Printf("GetCourseWithDetails: Fetched facilitator: %+v", facilitator)

	var jsonBooksResult json.RawMessage
	err3 := s.sb.DB.From("course_books").
		Select("book_id").
		Eq("course_id", courseID).
		Execute(&jsonBooksResult)
	if err3 != nil {
		log.Printf("GetCourseWithDetails: Error querying course books: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var courseBooks []models.CourseBook
	err = json.Unmarshal(jsonBooksResult, &courseBooks)
	if err != nil {
		log.Printf("GetCourseWithDetails: Error unmarshaling course books data: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}
	log.Printf("GetCourseWithDetails: Fetched course books: %+v", courseBooks)

	var books []models.Book
	for _, courseBook := range courseBooks {
		bookID := strconv.Itoa(courseBook.BookID)
		log.Printf("GetCourseWithDetails: Processing book ID: %s", bookID)

		var jsonBookResult json.RawMessage
		err := s.Sb().DB.From("books").
			Select("*").
			Single().
			Eq("id", bookID).
			Execute(&jsonBookResult)
		if err != nil {
			log.Printf("GetCourseWithDetails: Error querying book: %v", err)
			log.Printf("GetCourseWithDetails: Book ID: %s", bookID)
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
		}

		var book models.Book
		err = json.Unmarshal(jsonBookResult, &book)
		if err != nil {
			log.Printf("GetCourseWithDetails: Error unmarshaling book data: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
		}
		log.Printf("GetCourseWithDetails: Fetched book: %+v", book)

		books = append(books, book)
	}

	response := GetCourseWithDetailsResponse{
		Course:      course,
		Facilitator: facilitator,
		Books:       books,
	}
	log.Printf("GetCourseWithDetails: Response: %+v", response)

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

	var jsonResult json.RawMessage
	err = s.sb.DB.From("courses").Insert(string(data)).Execute(&jsonResult)
	if err != nil {
		log.Printf("Error inserting course: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Created course: %+v", course)
	return c.JSON(data)
}

func (s *Server) listFacilitators(c *fiber.Ctx) error {
	var jsonResult json.RawMessage
	err := s.sb.DB.From("facilitators").Select("*").Execute(&jsonResult)
	if err != nil {
		log.Printf("Error querying facilitators: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var facilitators []models.Facilitator
	err = json.Unmarshal(jsonResult, &facilitators)
	if err != nil {
		log.Printf("Error unmarshaling facilitators: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Facilitators: %+v", facilitators)
	return c.JSON(facilitators)
}

func (s *Server) getFacilitator(c *fiber.Ctx) error {
	facilitatorID := c.Params("id")

	var jsonResult json.RawMessage
	err := s.sb.DB.From("facilitators").
		Select("*").
		Eq("id", facilitatorID).
		Execute(&jsonResult)
	if err != nil {
		log.Printf("Error querying facilitator: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var facilitator models.Facilitator
	err = json.Unmarshal(jsonResult, &facilitator)
	if err != nil {
		log.Printf("Error unmarshaling facilitator: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Facilitator: %+v", facilitator)
	return c.JSON(facilitator)
}

func (s *Server) createFacilitator(c *fiber.Ctx) error {
	var facilitator models.Facilitator
	if err := c.BodyParser(&facilitator); err != nil {
		log.Printf("Error parsing facilitator: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: err.Error()})
	}

	data, err := json.Marshal(facilitator)
	if err != nil {
		log.Printf("Error marshaling facilitator: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var jsonResult json.RawMessage
	err = s.sb.DB.From("facilitators").Insert(string(data)).Execute(&jsonResult)
	if err != nil {
		log.Printf("Error inserting facilitator: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Created facilitator: %+v", facilitator)
	return c.JSON(facilitator)
}

func (s *Server) deleteFacilitator(c *fiber.Ctx) error {
	facilitatorID := c.Params("id")

	var jsonResult json.RawMessage
	err := s.sb.DB.From("facilitators").Delete().Eq("id", facilitatorID).Execute(&jsonResult)
	if err != nil {
		log.Printf("Error deleting facilitator: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Deleted facilitator with ID: %s", facilitatorID)
	return c.SendStatus(fiber.StatusNoContent)
}

// DISCUSSIONS
func (s *Server) createDiscussion(c *fiber.Ctx) error {
	var discussion models.Discussion
	if err := c.BodyParser(&discussion); err != nil {
		log.Printf("Error parsing discussion: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: err.Error()})
	}

	data, err := json.Marshal(discussion)
	if err != nil {
		log.Printf("Error marshaling discussion: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var jsonResult json.RawMessage
	err = s.sb.DB.From("discussions").Insert(string(data)).Execute(&jsonResult)
	if err != nil {
		log.Printf("Error inserting discussion: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Created discussion: %+v", discussion)
	return c.JSON(discussion)
}

func (s *Server) getDiscussion(c *fiber.Ctx) error {
	discussionID := c.Params("id")

	var jsonResult json.RawMessage
	err := s.sb.DB.From("discussions").
		Select("*").
		Single().
		Eq("id", discussionID).
		Execute(&jsonResult)
	if err != nil {
		log.Printf("Error querying discussion: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var discussion models.Discussion
	err = json.Unmarshal(jsonResult, &discussion)
	if err != nil {
		log.Printf("Error unmarshaling discussion: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Discussion: %+v", discussion)
	return c.JSON(discussion)
}

func (s *Server) listDiscussions(c *fiber.Ctx) error {
	var jsonResult json.RawMessage
	err := s.sb.DB.From("discussions").Select("*").Execute(&jsonResult)
	if err != nil {
		log.Printf("Error querying discussions: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var discussions []models.Discussion
	err = json.Unmarshal(jsonResult, &discussions)
	if err != nil {
		log.Printf("Error unmarshaling discussions: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Discussions: %+v", discussions)
	return c.JSON(discussions)
}

func (s *Server) updateDiscussion(c *fiber.Ctx) error {
	discussionID := c.Params("id")

	var discussion models.Discussion
	if err := c.BodyParser(&discussion); err != nil {
		log.Printf("Error parsing discussion: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: err.Error()})
	}

	data, err := json.Marshal(discussion)
	if err != nil {
		log.Printf("Error marshaling discussion: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var jsonResult json.RawMessage
	err = s.sb.DB.From("discussions").Update(string(data)).Eq("id", discussionID).Execute(&jsonResult)
	if err != nil {
		log.Printf("Error updating discussion: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Updated discussion: %+v", discussion)
	return c.JSON(discussion)
}

func (s *Server) deleteDiscussion(c *fiber.Ctx) error {
	discussionID := c.Params("id")

	var jsonResult json.RawMessage
	err := s.sb.DB.From("discussions").Delete().Eq("id", discussionID).Execute(&jsonResult)
	if err != nil {
		log.Printf("Error deleting discussion: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Deleted discussion with ID: %s", discussionID)
	return c.SendStatus(fiber.StatusNoContent)
}

//Readings!

func (s *Server) createReading(c *fiber.Ctx) error {
	var reading models.Reading
	if err := c.BodyParser(&reading); err != nil {
		log.Printf("Error parsing reading: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: err.Error()})
	}

	data, err := json.Marshal(reading)
	if err != nil {
		log.Printf("Error marshaling reading: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var jsonResult json.RawMessage
	err = s.sb.DB.From("readings").Insert(string(data)).Execute(&jsonResult)
	if err != nil {
		log.Printf("Error inserting reading: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Created reading: %+v", reading)
	return c.JSON(reading)
}

func (s *Server) getReading(c *fiber.Ctx) error {
	readingID := c.Params("id")

	var jsonResult json.RawMessage
	err := s.sb.DB.From("readings").
		Select("*").
		Single().
		Eq("id", readingID).
		Execute(&jsonResult)
	if err != nil {
		log.Printf("Error querying reading: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var reading models.Reading
	err = json.Unmarshal(jsonResult, &reading)
	if err != nil {
		log.Printf("Error unmarshaling reading: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Reading: %+v", reading)
	return c.JSON(reading)
}

func (s *Server) listReadings(c *fiber.Ctx) error {
	var jsonResult json.RawMessage
	err := s.sb.DB.From("readings").Select("*").Execute(&jsonResult)
	if err != nil {
		log.Printf("Error querying readings: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var readings []models.Reading
	err = json.Unmarshal(jsonResult, &readings)
	if err != nil {
		log.Printf("Error unmarshaling readings: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Readings: %+v", readings)
	return c.JSON(readings)
}

func (s *Server) updateReading(c *fiber.Ctx) error {
	readingID := c.Params("id")

	var reading models.Reading
	if err := c.BodyParser(&reading); err != nil {
		log.Printf("Error parsing reading: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: err.Error()})
	}

	data, err := json.Marshal(reading)
	if err != nil {
		log.Printf("Error marshaling reading: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var jsonResult json.RawMessage
	err = s.sb.DB.From("readings").Update(string(data)).Eq("id", readingID).Execute(&jsonResult)
	if err != nil {
		log.Printf("Error updating reading: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Updated reading: %+v", reading)
	return c.JSON(reading)
}

func (s *Server) deleteReading(c *fiber.Ctx) error {
	readingID := c.Params("id")

	var jsonResult json.RawMessage
	err := s.sb.DB.From("readings").Delete().Eq("id", readingID).Execute(&jsonResult)
	if err != nil {
		log.Printf("Error deleting reading: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Deleted reading with ID: %s", readingID)
	return c.SendStatus(fiber.StatusNoContent)
}

//Discussion Attendance

func (s *Server) createDiscussionAttendance(c *fiber.Ctx) error {
	var attendance models.DiscussionAttendance
	if err := c.BodyParser(&attendance); err != nil {
		log.Printf("Error parsing discussion attendance: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: err.Error()})
	}

	data, err := json.Marshal(attendance)
	if err != nil {
		log.Printf("Error marshaling discussion attendance: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var jsonResult json.RawMessage
	err = s.sb.DB.From("discussion_attendance").Insert(string(data)).Execute(&jsonResult)
	if err != nil {
		log.Printf("Error inserting discussion attendance: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Created discussion attendance: %+v", attendance)
	return c.JSON(attendance)
}

func (s *Server) listDiscussionAttendance(c *fiber.Ctx) error {
	discussionID := c.Params("id")

	var jsonResult json.RawMessage
	err := s.sb.DB.From("discussion_attendance").
		Select("*").
		Eq("discussion_id", discussionID).
		Execute(&jsonResult)
	if err != nil {
		log.Printf("Error querying discussion attendance: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var attendanceList []models.DiscussionAttendance
	err = json.Unmarshal(jsonResult, &attendanceList)
	if err != nil {
		log.Printf("Error unmarshaling discussion attendance: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Discussion attendance: %+v", attendanceList)
	return c.JSON(attendanceList)
}

// ReadingRating

func (s *Server) getReadingRating(c *fiber.Ctx) error {
	ratingID := c.Params("id")

	var jsonResult json.RawMessage
	err := s.sb.DB.From("reading_ratings").
		Select("*").
		Single().
		Eq("id", ratingID).
		Execute(&jsonResult)
	if err != nil {
		log.Printf("Error querying reading rating: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var rating models.ReadingRating
	err = json.Unmarshal(jsonResult, &rating)
	if err != nil {
		log.Printf("Error unmarshaling reading rating: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Reading rating: %+v", rating)
	return c.JSON(rating)
}

func (s *Server) listReadingRatings(c *fiber.Ctx) error {
	readingID := c.Params("id")

	var jsonResult json.RawMessage
	err := s.sb.DB.From("reading_ratings").
		Select("*").
		Eq("reading_id", readingID).
		Execute(&jsonResult)
	if err != nil {
		log.Printf("Error querying reading ratings: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var ratings []models.ReadingRating
	err = json.Unmarshal(jsonResult, &ratings)
	if err != nil {
		log.Printf("Error unmarshaling reading ratings: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Reading ratings: %+v", ratings)
	return c.JSON(ratings)
}

func (s *Server) createReadingRating(c *fiber.Ctx) error {
	var rating models.ReadingRating
	if err := c.BodyParser(&rating); err != nil {
		log.Printf("Error parsing reading rating: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: err.Error()})
	}

	data, err := json.Marshal(rating)
	if err != nil {
		log.Printf("Error marshaling reading rating: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var jsonResult json.RawMessage
	err = s.sb.DB.From("reading_ratings").Insert(string(data)).Execute(&jsonResult)
	if err != nil {
		log.Printf("Error inserting reading rating: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Created reading rating: %+v", rating)
	return c.JSON(rating)
}

func (s *Server) updateReadingRating(c *fiber.Ctx) error {
	ratingID := c.Params("id")

	var rating models.ReadingRating
	if err := c.BodyParser(&rating); err != nil {
		log.Printf("Error parsing reading rating: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: err.Error()})
	}

	data, err := json.Marshal(rating)
	if err != nil {
		log.Printf("Error marshaling reading rating: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var jsonResult json.RawMessage
	err = s.sb.DB.From("reading_ratings").Update(string(data)).Eq("id", ratingID).Execute(&jsonResult)
	if err != nil {
		log.Printf("Error updating reading rating: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Updated reading rating: %+v", rating)
	return c.JSON(rating)
}

func (s *Server) deleteReadingRating(c *fiber.Ctx) error {
	ratingID := c.Params("id")

	var jsonResult json.RawMessage
	err := s.sb.DB.From("reading_ratings").Delete().Eq("id", ratingID).Execute(&jsonResult)
	if err != nil {
		log.Printf("Error deleting reading rating: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	log.Printf("Deleted reading rating with ID: %s", ratingID)
	return c.SendStatus(fiber.StatusNoContent)
}
