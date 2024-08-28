package handlers

import (
	"context"

	messageservice "github.com/ssssunat/rest-api1/internal/messageService"
	messages "github.com/ssssunat/rest-api1/internal/web/messages"
)

type Handler struct {
	Service *messageservice.MessageService
}

func (h *Handler) GetMessages(_ context.Context, _ messages.GetMessagesRequestObject) (messages.GetMessagesResponseObject, error) {
	// Получение всех сообщений из сервиса
	allMessages, err := h.Service.GetAllMessages()
	if err != nil {
		return nil, err
	}

	// Создаем переменную респон типа 200джейсонРеспонс
	// Которую мы потом передадим в качестве ответа
	response := messages.GetMessages200JSONResponse{}

	// Заполняем слайс response всеми сообщениями из БД
	for _, msg := range allMessages {
		message := messages.Message{
			Id:      &msg.ID,
			Message: &msg.Text,
		}
		response = append(response, message)
	}

	// САМОЕ ПРЕКРАСНОЕ. Возвращаем просто респонс и nil!
	return response, nil
}

func (h *Handler) PostMessages(_ context.Context, request messages.PostMessagesRequestObject) (messages.PostMessagesResponseObject, error) {
	// Распаковываем тело запроса напрямую, без декодера!
	messageRequest := request.Body
	// Обращаемся к сервису и создаем сообщение
	messageToCreate := messageservice.Message{Text: *messageRequest.Message}
	createdMessage, err := h.Service.CreateMessage(messageToCreate)

	if err != nil {
		return nil, err
	}
	// создаем структуру респонс
	response := messages.PostMessages201JSONResponse{
		Id:      &createdMessage.ID,
		Message: &createdMessage.Text,
	}
	// Просто возвращаем респонс!
	return response, nil
}

func (h *Handler) DeleteMessagesId(ctx context.Context, request messages.DeleteMessagesIdRequestObject) (messages.DeleteMessagesIdResponseObject, error) {
	id := request.Id
	err := h.Service.DeleteMessageByID(int(id))
	if err != nil {
		return nil, err
	}
	return messages.DeleteMessagesId204Response{}, nil
}

func (h *Handler) PatchMessagesId(ctx context.Context, request messages.PatchMessagesIdRequestObject) (messages.PatchMessagesIdResponseObject, error) {
	id := request.Id
	messageReq := request.Body
	messageToUpdate := messageservice.Message{Text: *messageReq.Message}

	message, err := h.Service.UpdateMessageByID(int(id), messageToUpdate)
	if err != nil {
		return nil, err
	}
	response := messages.PatchMessagesId200JSONResponse{
		Id:      &message.ID,
		Message: &message.Text,
	}

	return response, nil
}

func NewHandler(service *messageservice.MessageService) *Handler {
	return &Handler{
		Service: service,
	}
}
