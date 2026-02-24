package res

type MergeStatus int

const (
	MergeWaiting MergeStatus = iota
	MergeMerging
	MergeSuccess
	MergeFailed
)

type MergeState struct {
	Status MergeStatus `json:"status"`
	Merged int         `json:"merged"`
	Total  int         `json:"total"`
	Error  string      `json:"error"`
}
