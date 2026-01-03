package mq_eneity

type TranscodeTask struct {
	UploadID string `json:"uploadId"`
	TmepFile string `json:"tempFile"`
	OutFile  string `json:"outFile"`
	RoomID   string `json:"roomId"`
}
