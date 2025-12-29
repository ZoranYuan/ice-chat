package repository

import (
	"errors"
	model "ice-chat/internal/model/eneity"
	"ice-chat/pkg/mysql"

	"gorm.io/gorm"
)

type groupsRepoImpl struct {
	db *mysql.DBUtils
}

type RoomsRepository interface {
	Create(group *model.Rooms, groupMembers *model.RoomsMember) error
	GroupIsExists(groupId uint64) (bool, error)
}

func NewRoomsRepo(db *mysql.DBUtils) RoomsRepository {
	return &groupsRepoImpl{
		db: db,
	}
}

func (r *groupsRepoImpl) Create(room *model.Rooms, roomsMembers *model.RoomsMember) error {
	err := r.db.Client().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&room).Error; err != nil {
			return err
		}

		roomsMembers.RoomID = room.RoomID

		if err := tx.Create(&roomsMembers).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.New("room already exists")
		}
		return err
	}

	return err
}

func (r *groupsRepoImpl) GroupIsExists(roomId uint64) (bool, error) {
	var count int64
	err := r.db.Client().Model(&model.Rooms{}).Where("room_id = ?", roomId).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count == 1, nil
}
