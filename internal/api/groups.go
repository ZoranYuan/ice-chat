package api

import (
	"ice-chat/internal/model/request"
	"ice-chat/internal/response"
	"ice-chat/internal/service"

	"github.com/gin-gonic/gin"
)

type GroupApi interface {
	Create(ctx *gin.Context)
}

type groupApiImpl struct {
	groupServ service.GroupsService
}

func NewGroupsApi(groupServ service.GroupsService) GroupApi {
	return &groupApiImpl{
		groupServ: groupServ,
	}
}

func (g groupApiImpl) Create(ctx *gin.Context) {
	var group request.Group
	if err := ctx.ShouldBindBodyWithJSON(&group); err != nil {
		response.BadRequest(ctx)
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

	err, groupId := g.groupServ.Create(group, uid)
	if err != nil {
		response.Fail(ctx, 201, err.Error())
		ctx.Abort()
		return
	}

	response.OKWithData(ctx, gin.H{
		`groupId`: groupId,
	})
}
