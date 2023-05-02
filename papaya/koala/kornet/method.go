package kornet

const (
  GET = iota
  HEAD
  POST
  PUT
  DELETE
  CONNECT
  OPTIONS
  TRACE
)

func KHttpGetMethod(value int) string {

  switch value {
  case GET:
    return "GET"
  case HEAD:
    return "HEAD"
  case POST:
    return "POST"
  case PUT:
    return "PUT"
  case DELETE:
    return "DELETE"
  case CONNECT:
    return "CONNECT"
  case OPTIONS:
    return "OPTIONS"
  case TRACE:
    return "TRACE"
  }

  return "GET"
}
