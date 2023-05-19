package kornet

import (
  "math/rand"
  "time"
)

type Result struct {
  ID      uint     `json:"id"`
  Logs    []string `json:"logs"`
  Status  string   `json:"status"`
  Message string   `json:"message"`
  Error   bool     `json:"error"`
  Data    any      `json:"data"`
}

func ResultNew(message *Message, data any) *Result {

  // fake result ID, just wrapping message and data
  randomize := rand.New(rand.NewSource(time.Now().UnixNano()))

  min := 1000000000
  max := 2147483647

  return &Result{
    ID:      uint(randomize.Intn(max-min+1) + min), // randomized
    Logs:    make([]string, 0),
    Status:  "Unknown", // Modify by Context
    Message: message.Message,
    Error:   message.Error,
    Data:    data,
  }
}
