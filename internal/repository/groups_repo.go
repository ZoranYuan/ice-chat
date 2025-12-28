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

type GroupsRepository interface {
	Create(group *model.Groups, groupMembers *model.GroupMember) error
}

func NewGroupsRepo(db *mysql.DBUtils) GroupsRepository {
	return &groupsRepoImpl{
		db: db,
	}
}

func (g *groupsRepoImpl) Create(group *model.Groups, groupMembers *model.GroupMember) error {
	err := g.db.Client().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&group).Error; err != nil {
			return err
		}

		groupMembers.GroupID = group.GroupId

		if err := tx.Create(&groupMembers).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.New("group already exists")
		}
		return err
	}

	return err
}
