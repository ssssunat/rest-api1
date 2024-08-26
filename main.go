package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)
type Message struct {
	Message string `json:"message"`
}
var m Message
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, " + m.Message)
}

func apiMessage(w http.ResponseWriter, r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		fmt.Fprintf(w, "error json Body")
		return
	} 
	json.NewEncoder(w).Encode(m)
}

func main() {
	router := chi.NewRouter()
	router.Get("/api/hello", HelloHandler)
	router.Post("/api/message", apiMessage)
	http.ListenAndServe(":8080", router)
}
