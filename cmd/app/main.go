package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ssssunat/rest-api1/internal/database"
	"github.com/ssssunat/rest-api1/internal/handlers"
	messageservice "github.com/ssssunat/rest-api1/internal/messageService"
	messages "github.com/ssssunat/rest-api1/internal/web/messages"
)

func main() {
	database.InitDB()
	if err := database.DB.AutoMigrate(&messageservice.Message{}); err != nil {
		log.Fatalf("Failed AutoMigrate DB: %v", err)
	}
	repo := messageservice.NewMessageRepository(database.DB)
	service := messageservice.NewService(repo)

	handler := handlers.NewHandler(service)
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	strictHandler := messages.NewStrictHandler(handler, nil)
	messages.RegisterHandlers(e, strictHandler)

	if err := e.Start(":8080"); err != nil {
		log.Fatalf("Failed to start with err: %v", err)
	}
}
