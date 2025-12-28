package repository

import (
	model "ice-chat/internal/model/eneity"
	"ice-chat/pkg/mysql"
	"log"
)

type userRepo struct {
	db *mysql.DBUtils
}

type UserRepository interface {
	FindUserByEmail(email string) (*model.User, error)
	IsUserExist(id uint64) bool
}

func NewUserRepository(db *mysql.DBUtils) UserRepository {
	return &userRepo{db: db}
}

func (u *userRepo) FindUserByEmail(email string) (*model.User, error) {
	var user *model.User
	err := u.db.Client().Model(&model.User{}).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userRepo) IsUserExist(id uint64) bool {
	if id == 0 {
		return false
	}

	var count int64
	err := u.db.Client().Model(&model.User{}).Where("id = ?", id).Count(&count).Error

	if err != nil {
		log.Printf("failed to check user by : %d, error: %v", id, err)
		return false
	}

	return count > 0
}
