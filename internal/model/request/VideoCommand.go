package request

type VideoAction string

const (
	ActionPlay       VideoAction = "play"
	ActionPause                  = "pause"
	ActionSeek                   = "seek"
	ActionChangeRate             = "change_rate"
)

type VideoCommand struct {
	RoomId    uint64      `json:"room_id" binding:"required"`
	Action    VideoAction `json:"action" binding:"required"`
	Time      float64     `json:"time,omitempty"` // seek 时使用，前端拖拽进度条的时候使用
	Speed     float64     `json:"speed,omitempty"`
	TimeStamp int64       `json:"time_stamp"` // 防止旧消息
}
