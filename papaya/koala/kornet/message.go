package kornet

type Message struct {
  Message string `json:"message"`
  Error   bool   `json:"error"`
}

func MessageNew(message string, bad bool) *Message {

  return &Message{
    Message: message,
    Error:   bad,
  }
}

// shorthand

func Msg(message string, bad bool) *Result {

  return ResultNew(MessageNew(message, bad), nil)
}
