package kornet

import (
	"math/rand"
)

type Result struct {
	ID      uint     `json:"id"`
	Logs    []string `json:"logs"`
	Status  string   `json:"status"`
	Message string   `json:"string"`
	Error   bool     `json:"error"`
	Data    any      `json:"data"`
}

func ResultNew(message *Message, data any) *Result {

	// fake result ID, just wrapping message and data

	return &Result{
		ID:      uint(rand.Uint32()), // randomized
		Logs:    make([]string, 0),
		Status:  "Unknown", // Modify by Context
		Message: message.Message,
		Error:   message.Error,
		Data:    data,
	}
}

func MessageOnly(message *Message) *Result {

	// fake result ID, just wrapping message only

	return &Result{
		ID:      uint(rand.Uint32()), // randomized
		Logs:    make([]string, 0),
		Status:  "Unknown", // Modify by Context
		Message: message.Message,
		Error:   message.Error,
		Data:    nil,
	}
}
