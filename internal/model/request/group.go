package request

type Group struct {
	GroupName string `json:"groupName" binding:"required"`
	Avatar    string `json:"avatar"`
	Desc      string `json:"desc"`
}
