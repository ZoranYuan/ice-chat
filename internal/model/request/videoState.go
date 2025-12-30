package request

type VideoState struct {
	RoomID    uint64  `json:"roomId" binding:"required"`
	Progress  float64 `json:"progress" binding:"required"`
	IsPlaying bool    `json:"isPlaying" binding:"required"`
	Speed     float32 `json:"speed" binding:"required"`
	TimeStamp uint64  `json:"timeStamp" binding:"required"`
}
