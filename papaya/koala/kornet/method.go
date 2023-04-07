package kornet

const (
  MethodGet = iota
  MethodHead
  MethodPost
  MethodPut
  MethodDelete
  MethodConnect
  MethodOptions
  MethodTrace
)

func KHttpGetMethod(value int) string {

  switch value {
  case MethodGet:
    return "GET"
  case MethodHead:
    return "HEAD"
  case MethodPost:
    return "POST"
  case MethodPut:
    return "PUT"
  case MethodDelete:
    return "DELETE"
  case MethodConnect:
    return "CONNECT"
  case MethodOptions:
    return "OPTIONS"
  case MethodTrace:
    return "TRACE"
  }

  return "GET"
}
