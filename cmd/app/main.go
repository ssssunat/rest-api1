package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ssssunat/rest-api1/internal/database"
	"github.com/ssssunat/rest-api1/internal/handlers"
	messageservice "github.com/ssssunat/rest-api1/internal/messageService"
)




func main() {
	database.InitDB()
	database.DB.AutoMigrate(&messageservice.Message{})

	repo := messageservice.NewMessageRepository(database.DB)
	service := messageservice.NewService(repo)

	handler := handlers.NewHandler(service)

	router := chi.NewRouter()

	router.Get("/api/get", handler.GetMessagesHandler)
	router.Patch("/api/update/{id}", handler.UpdateMessageHandler)
	router.Delete("/api/delete/{id}", handler.DeleteMessageHandler)
	router.Post("/api/post", handler.PostMessageHandler)

	http.ListenAndServe(":8080", router)
}
