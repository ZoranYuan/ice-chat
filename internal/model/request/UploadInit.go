package request

type UploadInit struct {
	FileName string `json:"fileName" binding:"required"`
	Size     int64  `json:"size" binding:"required"`
}
