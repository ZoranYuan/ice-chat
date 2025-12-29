package request

type UploadMerge struct {
	UploadId    string `json:"uploadId" binding:"required"`
	TotalChunks int    `json:"totalChunks" binding:"required"`
}
