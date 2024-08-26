package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=localhost user=postgres password=yourpassword dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connet to database:", err)
	}
	DB = db
}

type Message struct {
	gorm.Model
	Text string `json:"text"` // Наш сервер будет ожидать json c полем text
}

type requestBody struct {
	Message string `json:"message"`
}

var m requestBody

func GetMessage(w http.ResponseWriter, r *http.Request) {
	var messages []Message
	DB.Find(&messages)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Invalid input", http.StatusNotFound)
		return
	}

	message := Message{Text: m.Message}
	DB.Create(&message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(m)
}

func UpdateMessage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var message Message
	if err := DB.First(&message, id).Error; err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	var input Message
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	DB.Model(&message).Updates(input)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	log.Println("DeleteMessage handler called")
	id := chi.URLParam(r, "id")
	if err := DB.Delete(&Message{}, id).Error; err != nil {
		log.Println("Error deleting message:", err)
		http.Error(w, "Message ot found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	InitDB()

	DB.AutoMigrate(&Message{})
	router := chi.NewRouter()

	router.Get("/api/messages", GetMessage)
	router.Patch("/api/messages/{id}", UpdateMessage)
	router.Delete("/api/messages/{id}", DeleteMessage)
	router.Post("/api/messages", CreateMessage)

	http.ListenAndServe(":8080", router)
}
