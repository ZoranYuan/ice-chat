package request

type UploadMerge struct {
	RoomId      uint64 `json:"roomId" binding:"required"`
	UploadId    string `json:"uploadId" binding:"required"`
	TotalChunks int    `json:"totalChunks" binding:"required"`
}
