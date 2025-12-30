package api

import (
	"ice-chat/internal/model/request"
	res "ice-chat/internal/model/response"
	"ice-chat/internal/response"
	"ice-chat/internal/service"

	"github.com/gin-gonic/gin"
)

type RoomsApi interface {
	CreateRoom(c *gin.Context)
	JoinRoom(c *gin.Context)
}

type roomsApiImpl struct {
	roomServ service.RoomsService
}

func NewRoomsApi(groupServ service.RoomsService) RoomsApi {
	return &roomsApiImpl{
		roomServ: groupServ,
	}
}

func (r *roomsApiImpl) CreateRoom(c *gin.Context) {
	var room request.Room
	if err := c.ShouldBindBodyWithJSON(&room); err != nil {
		response.BadRequestWithMessage(c, "参数错误")
		c.Abort()
		return
	}

	v, exists := c.Get("uid")
	uid, ok := v.(uint64)
	if !exists || !ok {
		response.Unauthorized(c)
		c.Abort()
		return
	}

	err, roomId, joinCode := r.roomServ.CreateRoom(room, uid)
	if err != nil {
		response.Fail(c, 201, err.Error())
		c.Abort()
		return
	}

	response.OKWithData(c, res.Room{
		RoomID:   roomId,
		JoinCode: joinCode,
	})
}

func (r *roomsApiImpl) JoinRoom(c *gin.Context) {
	joinCode := c.Param("joinCode")
	if joinCode == "" || len(joinCode) != 6 {
		response.BadRequestWithMessage(c, "无效的验证码")
		c.Abort()
		return
	}

	v, exists := c.Get("uid")
	uid, ok := v.(uint64)
	if !exists || !ok {
		response.Unauthorized(c)
		c.Abort()
		return
	}

	roomId, err := r.roomServ.JoinRoom(uid, joinCode)
	if err != nil {
		response.Fail(c, 201, err.Error())
		c.Abort()
		return
	}

	response.OKWithData(c, gin.H{
		"roomId": roomId,
	})
}
