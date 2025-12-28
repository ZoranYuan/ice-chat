package request

type Room struct {
	RoomName string `json:"roomName" binding:"required"`
	Avatar   string `json:"avatar"`
	Desc     string `json:"desc"`
}
