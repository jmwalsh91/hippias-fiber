package server

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/supabase-community/supabase-go"
)

type Server struct {
	*fiber.App
	Sb *supabase.Client
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
		Sb:  client,
	}
	server.App.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	return server
}
