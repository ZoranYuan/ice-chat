package api

import (
	"fmt"
	"ice-chat/internal/constants"
	cerror "ice-chat/internal/model/error"
	"ice-chat/internal/model/request"
	res "ice-chat/internal/model/response"
	"ice-chat/internal/response"
	"ice-chat/internal/service"
	"ice-chat/utils"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadApi interface {
	UploadInit(c *gin.Context)
	Upload(c *gin.Context)
	Merge(c *gin.Context)
}

type uploadApiImpl struct {
	uploadSev service.UploadService
}

func NewUploadApi(uploadSev service.UploadService) UploadApi {
	return &uploadApiImpl{
		uploadSev: uploadSev,
	}
}

func (u *uploadApiImpl) UploadInit(c *gin.Context) {
	// TODO 为当前要上传的文件生成 uuid，并且指定分块原则
	var fileInfo request.UploadInit
	if err := c.ShouldBindBodyWithJSON(&fileInfo); err != nil {
		response.BadRequestWithMessage(c, err.Error())
		c.Abort()
		return
	}

	_, isValid := utils.IsValidFileExt(fileInfo.FileName)
	if !isValid {
		response.BadRequestWithMessage(c, "视频格式错误")
		c.Abort()
		return
	}

	// TODO 根据 file 的 size 来决定使用什么分片方式，这里统一
	uploadId := fmt.Sprintf("%s_%d", uuid.New().String(), time.Now().Unix())
	response.OKWithData(c, res.UploadInit{
		UploadId:        uploadId,
		UploadChunkSize: int64(constants.UPLOAD_CHUNK_SIZE),
	})
}

func (u *uploadApiImpl) Upload(c *gin.Context) {
	uploadId := c.PostForm("uploadId")
	uploadChunkIdxstr := c.PostForm("uploadChunkIdx")

	if uploadId == "" || uploadChunkIdxstr == "" {
		response.BadRequestWithMessage(c, "参数错误")
		c.Abort()
		return
	}

	uploadChunkIdx, _ := strconv.Atoi(uploadChunkIdxstr)

	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequestWithMessage(c, "参数错误")
		c.Abort()
		return
	}

	dir := fmt.Sprintf("%s/%s", constants.UPLOAD_TEMP_DIR, uploadId)
	os.MkdirAll(dir, 0755)

	dest := fmt.Sprintf("%s/%d", dir, uploadChunkIdx)

	if err := c.SaveUploadedFile(file, dest); err != nil {
		response.InternalError(c, "网络错误，请稍后重试")
		c.Abort()
		return
	}

	if err := u.uploadSev.Upload(uploadId, uploadChunkIdx); err != nil {
		response.InternalError(c, "获取数据失败，请稍后重试")
		c.Abort()
		return
	}

	response.OK(c)
}

func (u *uploadApiImpl) Merge(c *gin.Context) {
	var mergeInfo request.UploadMerge
	if err := c.ShouldBindBodyWithJSON(&mergeInfo); err != nil {
		response.BadRequestWithMessage(c, "参数错误")
		c.Abort()
		return
	}

	err := u.uploadSev.Merge(mergeInfo)

	if err != nil {
		if e, ok := err.(*cerror.ChunkMissingError); ok {
			response.BadRequestWithData(c, res.UploadMerge{
				UploadId:     mergeInfo.UploadId,
				IsLost:       true,
				LostChunkIdx: e.MissingIndex,
			})
		}
	}

	response.OK(c)
}

func (u *uploadApiImpl) BreakPointContinue(c *gin.Context) {
	// TODO 获取断点
	uploadId := c.Param("uploadId")

	// TODO 检测这个 ID 是否存在，后续完成
	current, err := u.uploadSev.BreakPoint(uploadId)
	if err != nil {
		response.Fail(c, 201, "获取失败，请重新上传")
		c.Abort()
		return
	}

	response.OKWithData(c, res.UploadPoint{
		UploadChunkIdx: int(current),
	})
}
