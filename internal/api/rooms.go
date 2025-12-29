package api

import (
	"ice-chat/internal/model/request"
	res "ice-chat/internal/model/response"
	"ice-chat/internal/response"
	"ice-chat/internal/service"

	"github.com/gin-gonic/gin"
)

type RoomsApi interface {
	Create(ctx *gin.Context)
}

type roomsApiImpl struct {
	groupServ service.RoomsService
}

func NewRoomsApi(groupServ service.RoomsService) RoomsApi {
	return roomsApiImpl{
		groupServ: groupServ,
	}
}

func (r roomsApiImpl) Create(ctx *gin.Context) {
	var group request.Room
	if err := ctx.ShouldBindBodyWithJSON(&group); err != nil {
		response.BadRequestWithMessage(ctx, "参数错误")
		ctx.Abort()
		return
	}

	v, exists := ctx.Get("uid")
	uid, ok := v.(uint64)
	if !exists || !ok {
		response.Unauthorized(ctx)
		ctx.Abort()
		return
	}

	err, roomId := r.groupServ.Create(group, uid)
	if err != nil {
		response.Fail(ctx, 201, err.Error())
		ctx.Abort()
		return
	}

	response.OKWithData(ctx, res.Room{
		RoomID: roomId,
	})
}
