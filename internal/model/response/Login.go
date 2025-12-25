package response

type Login struct {
	Token    string `json:"token"`
	UserId   uint64 `json:"userId"`
	UserName string `json:"userName"`
	Avatar   string `json:"avatar"`
}
