package repository

import (
	model "ice-chat/internal/model/eneity"
	"ice-chat/internal/model/request"
	"ice-chat/pkg/mysql"
)

type msgRepo struct {
	db *mysql.DBUtils
}

type MessageRepository interface {
	Add(message request.Message) error
}

func NewUmsgRepository(db *mysql.DBUtils) MessageRepository {
	return &msgRepo{
		db: db,
	}
}

func (m *msgRepo) Add(message request.Message) error {
	err := m.db.Client().Create(&model.Message{
		UserID:      message.UserId,
		Content:     message.Content,
		Username:    message.UserName,
		MessageType: message.MessageType,
		Status:      1,
		CreatedAt:   message.Time,
	}).Error
	return err
}
