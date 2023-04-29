package cors

import (
  "net/url"
  m "skfw/papaya/bunny/swag/method"
  "strconv"
  "strings"
)

func SafeURL(url *url.URL) string {

  scheme := "https"
  host := "localhost" // fallback into local network

  if url.Scheme != "" {

    scheme = url.Scheme
  }

  if url.Host != "" {

    host = url.Host
  }

  return scheme + "://" + host
}

func SafeMethods(methods []string) string {

  if len(methods) > 0 {

    temp := make([]string, 0)

    for _, method := range methods {

      switch method {

      case m.GET, m.HEAD, m.POST, m.PUT, m.DELETE, m.CONNECT, m.OPTIONS, m.TRACE:

        temp = append(temp, method)
      }
    }

    return strings.Join(temp, ",")

  }

  return "GET" // fallback into readonly
}

func SafeHeaders(headers []string) string {

  if len(headers) > 0 {

    return strings.Join(headers, ",")
  }

  return "Content-Type" // fallback into defult value
}

func SafeCredentials(credentials bool) string {

  // not safe anymore, cause couldn't look up Origin first
  return strconv.FormatBool(credentials)
}

func SafeMaxAge(maxAge int) string {

  if maxAge > 0 {

    return strconv.Itoa(maxAge)
  }

  return "86400"
}
