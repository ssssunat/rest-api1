package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	messageservice "github.com/ssssunat/rest-api1/internal/messageService"
)

type Handler struct {
	Service *messageservice.MessageService
}

func NewHandler(service *messageservice.MessageService) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	messages, err := h.Service.GetAllMessages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func (h *Handler) PostMessageHandler(w http.ResponseWriter, r *http.Request) {
	var message messageservice.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdMessage, err := h.Service.CreateMessage(message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdMessage)
}

func (h *Handler) DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Service.DeleteMessageByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateMessageHandler(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var message messageservice.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, err = h.Service.UpdateMessageByID(id, message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")  // это обезательно писать
	w.WriteHeader(http.StatusOK)										// 		
	json.NewEncoder(w).Encode(message) // нужно ли делать обработку ...
}
