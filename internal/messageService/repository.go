package messageservice

import "gorm.io/gorm"

type MessageRepository interface {
	CreateMessage(message Message) (Message, error)
	GetAllMessages() ([]Message, error)
	UpdateMessageByID(id int, message Message) (Message, error)
	DeleteMessageByID(id int) error
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *messageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) CreateMessage(message Message) (Message, error) {
	result := r.db.Create(&message)
	if result.Error != nil {
		return Message{}, result.Error
	}
	return message, nil
}

func (r *messageRepository) GetAllMessages() ([]Message, error) {
	var messages []Message
	err := r.db.Find(&messages).Error
	return messages, err
}

func (r *messageRepository) UpdateMessageByID(id int, message Message) (Message, error) {
	var messageUpdate Message
	err := r.db.First(&messageUpdate, id).Error
	if err != nil {
		return Message{}, err
	}

	err = r.db.Model(&messageUpdate).Updates(message).Error
	if err != nil {
		return Message{}, err
	}

	return messageUpdate,nil
}

func (r *messageRepository) DeleteMessageByID(id int) error {
	err := r.db.Delete(&Message{}, id).Error
	return err
}