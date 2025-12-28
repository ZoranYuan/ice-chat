package request

import "time"

type Message struct {
	Content     string    `json:"content" binding:"required"`
	UserId      uint64    `json:"userId" binding:"required"`
	UserName    string    `json:"userName" binding:"required"`
	MessageType int8      `json:"messageType"`
	Time        time.Time `json:"time"`
}
