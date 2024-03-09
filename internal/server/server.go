package server

import (
	"os"

	"github.com/gofiber/fiber/v2"
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
	server := &Server{
		App: fiber.New(),
		sb:  client,
	}

	server.App.Get("/book/id", server.getBook)
	server.App.Get("/list", server.listBooks)
	server.App.Get("/authors", server.listAuthors)
	server.App.Get("/authors/:id", server.getAuthor)
	server.App.Get("/books/author/:id", server.getBooksByAuthorID)

	return server
}
