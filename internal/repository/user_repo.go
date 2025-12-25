package repository

import (
	model "ice-chat/internal/model/eneity"
	"ice-chat/pkg/mysql"
)

type userRepo struct {
	db *mysql.DBUtils
}

type UserRepository interface {
	FindUserByEmail(email string) (*model.User, error)
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
