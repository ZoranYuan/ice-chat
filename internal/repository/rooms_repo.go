package repository

import (
	"errors"
	model "ice-chat/internal/model/eneity"
	my_mysql "ice-chat/pkg/mysql"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type groupsRepoImpl struct {
	db *my_mysql.DBUtils
}

type RoomsRepository interface {
	Create(group *model.Rooms, groupMembers *model.RoomsMember) error
	RoomIsExists(groupId uint64) (bool, error)
	JoinRoom(uid, roomId uint64) error
}

func NewRoomsRepo(db *my_mysql.DBUtils) RoomsRepository {
	return &groupsRepoImpl{
		db: db,
	}
}

// 通用判断函数
func IsDuplicateKeyErr(err error) bool {
	if err == nil {
		return false
	}
	// 匹配 GORM 封装的错误
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	// 匹配 MySQL 原生错误
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return true
	}
	return false
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

func (r *groupsRepoImpl) RoomIsExists(roomId uint64) (bool, error) {
	var count int64
	err := r.db.Client().Model(&model.Rooms{}).Where("room_id = ?", roomId).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count == 1, nil
}

func (r *groupsRepoImpl) JoinRoom(uid, roomId uint64) error {
	return r.db.Client().Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&model.Rooms{}).Where("room_id = ?", roomId).Count(&count).Error; err != nil {
			return err
		}

		if count != 1 {
			return errors.New("房间不存在")
		}

		var roomsMember = model.RoomsMember{
			UserID: uid,
			RoomID: roomId,
		}

		err := r.db.Client().Create(&roomsMember).Error
		// 幂等：已经是成员
		if IsDuplicateKeyErr(err) {
			return errors.New("你已经在房间中，请勿重复加入")
		}

		return nil
	})
}
