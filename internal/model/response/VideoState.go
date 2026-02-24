package res

type VideoState struct {
	RoomID      uint64  `json:"roomId" binding:"required"`
	CurrentTime float64 `json:"currentTime" binding:"required"`
	Duration    float64 `json:"duration" binding:"required"`
	IsPlaying   bool    `json:"isPlaying" binding:"required"`
	Speed       float64 `json:"speed" binding:"required"`
	UpdatedAt   int64   `json:"updatedAt"`
}
