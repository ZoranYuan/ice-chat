package res

type VideoStateInit struct {
	RoomID    string  `json:"roomId" binding:"required"`
	Progress  float64 `json:"progress" binding:"required"`
	IsPlaying bool    `json:"isPlaying" binding:"required"`
	Speed     float32 `json:"speed" binding:"required"`
	TimeStamp int64   `json:"timeStamp" binding:"required"`
	VideoUrl  string  `json:"videoUrl"`
}
