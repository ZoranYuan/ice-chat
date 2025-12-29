package res

type Message struct {
	Error       bool   `json:"error"`
	Content     string `json:"content" binding:"required"`
	UserId      uint64 `json:"userId" binding:"required"`
	UserName    string `json:"userName" binding:"required"`
	MessageType int8   `json:"messageType"`
	Time        string `json:"time"`
}
