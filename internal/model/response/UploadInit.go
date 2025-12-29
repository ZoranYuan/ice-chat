package res

type UploadInit struct {
	UploadId        string `json:"uploadId" binding:"required"`
	UploadChunkSize int64  `json:"uploadChunkSize" binding:"required"`
}
