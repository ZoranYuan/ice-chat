package api

import (
	"errors"
	"fmt"
	"ice-chat/internal/constants"
	"ice-chat/internal/model/request"
	"ice-chat/internal/response"
	"ice-chat/internal/service"
	"ice-chat/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (u *UserAPI) Upload(ctx *gin.Context) {
	userId := ctx.GetUint64("uid")
	if userId == 0 {
		response.Unauthorized(ctx)
		ctx.Abort()
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		response.BadRequestWithMessage(ctx, "参数错误")
		ctx.Abort()
		return
	}

	if file.Size > int64(constants.FILE_SIZE_1G)*5 {
		response.Fail(ctx, 201, "文件大小不能超过 5g")
		ctx.Abort()
		return
	}

	src, err := file.Open()
	if err != nil {
		response.InternalError(ctx, "文件读取失败，请稍后重试")
		ctx.Abort()
		return
	}
	defer src.Close()

	fileExt, isValid := utils.IsValidFileExt(file.Filename)
	if !isValid {
		response.Fail(ctx, 201, "仅支持上传.mp4/.avi/.mov/.mkv/.flv/.wmv格式的视频文件")
		ctx.Abort()
		return
	}

	objectKey := fmt.Sprintf("videos/%d/%s.%s", userId, uuid.New().String(), fileExt)
}
