package cors

import (
  m "skfw/papaya/bunny/swag/method"
  "strings"
)

func CheckMethodAvail(method string) bool {

  switch method {

  case m.GET, m.HEAD, m.POST, m.PUT, m.DELETE, m.CONNECT, m.OPTIONS, m.TRACE:

    return true
  }

  return false
}

func NormListBySources(sources []string, data []string) []string {

  var found bool
  var temp []string

  if len(data) > 0 {

    temp = make([]string, 0)

    for _, curr := range sources {

      curr = strings.ToUpper(curr)

      // find dup
      found = false
      for _, value := range temp {

        // case insensitive
        if value == curr {

          found = true
          break
        }
      }

      // no dup
      if !found {

        // check available
        for _, value := range data {

          // case insensitive
          if strings.ToUpper(value) == curr {

            temp = append(temp, curr)
            break
          }
        }
      }
    }

  } else {

    temp = sources
  }

  return temp
}
