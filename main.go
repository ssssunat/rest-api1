package main

import (
	"encoding/json"
	"fmt"
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
	result := DB.Find(&messages)
	if result.Error != nil {
		log.Printf("Error while fetch messages %s", result.Error)
	}
	json.NewEncoder(w).Encode(messages)
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		fmt.Fprintf(w, "error json Body")
		return
	}

	message := Message{Text: m.Message}
	result := DB.Create(&message)
	if result.Error != nil {
		log.Println(result.Error)
	}
	json.NewEncoder(w).Encode(m)
}

func main() {
	InitDB()

	DB.AutoMigrate(&Message{})
	router := chi.NewRouter()

	router.Get("/api/messages", GetMessage)
	router.Post("/api/messages", CreateMessage)

	http.ListenAndServe(":8080", router)
}
