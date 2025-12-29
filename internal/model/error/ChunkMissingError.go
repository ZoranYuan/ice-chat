package cerror

import "fmt"

type ChunkMissingError struct {
	MissingIndex int `json:"missingIndex"`
}

func (e *ChunkMissingError) Error() string {
	return fmt.Sprintf("chunk %d missing", e.MissingIndex)
}
