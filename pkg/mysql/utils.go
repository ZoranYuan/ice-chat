package my_mysql

import (
	model "ice-chat/internal/model/eneity"

	"gorm.io/gorm"
)

type DBUtils struct {
	client *gorm.DB
}

func NewDBUtils(client *gorm.DB) *DBUtils {
	return &DBUtils{client: client}
}

func (du *DBUtils) Client() *gorm.DB {
	return du.client
}

func (du *DBUtils) AutoMigrate() error {
	return du.client.AutoMigrate(&model.User{}, &model.Message{}, &model.Rooms{}, &model.RoomsMember{})
}
