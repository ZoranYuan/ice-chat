package utils

import (
	"bytes"
	"encoding/json"
	"log"
)

func UnmarshalTools(msg []byte, s any) bool {
	decoder := json.NewDecoder(bytes.NewReader(msg))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&s); err != nil {
		log.Printf("failed to marsha1: %v", err)
		return false
	}
	return true
}
