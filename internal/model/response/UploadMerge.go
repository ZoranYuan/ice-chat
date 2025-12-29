package res

type UploadMerge struct {
	UploadId     string `json:"uploadId"`
	IsLost       bool   `json:"isLost"`
	LostChunkIdx int    `json:"lostChunkIdx"`
}
