package api

import (
	"errors"
	"ice-chat/internal/constants"
	"ice-chat/internal/model/request"
	"ice-chat/internal/response"
	"ice-chat/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserAPI API层：处理用户相关HTTP请求
type UserAPI struct {
	userSvc *service.UserService // 依赖业务层
}

// NewUserAPI 构造函数：注入业务层实例
func NewUserAPI(userSvc *service.UserService) *UserAPI {
	return &UserAPI{
		userSvc: userSvc,
	}
}

func (u *UserAPI) Login(ctx *gin.Context) {
	var v request.Login
	if err := ctx.ShouldBindBodyWithJSON(&v); err != nil {
		response.BadRequestWithMessage(ctx, "参数错误")
		ctx.Abort()
		return
	}
	res, err := u.userSvc.Login(v)
	if err != nil {
		msg := ""
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg = constants.USERNOTFOUND
		} else {
			msg = err.Error()
		}
		response.Fail(ctx, 201, msg)
		ctx.Abort()
		return
	}

	response.OKWithData(ctx, res)
}
