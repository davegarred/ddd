package ddd

import (
	"encoding/json"
)


type ErrorDto struct {
	Error            string          `json:"error"`
	ValidationErrors *[]ErrorDetails `json:"validation_errors,omitempty"`
}

type ErrorDetails struct {
	Property string `json:"property"`
	Details  string `json:"details"`
}

func cleanlySerialize(dto interface{}) string {
	ser, err := json.Marshal(dto)
	if err != nil {
		panic(err)
	}
	return string(ser)
}
