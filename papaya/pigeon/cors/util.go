package cors

import (
  m "skfw/papaya/bunny/swag/method"
)

func CheckMethodAvail(method string) bool {

  switch method {

  case m.GET, m.HEAD, m.POST, m.PUT, m.DELETE, m.CONNECT, m.OPTIONS, m.TRACE:

    return true
  }

  return false
}
