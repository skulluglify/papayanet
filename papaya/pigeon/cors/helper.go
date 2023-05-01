package cors

import (
  "net/url"
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

    // continuously
    for _, vMethod := range Methods {

      vMethod = strings.ToUpper(vMethod)

      for _, cMethod := range methods {

        if strings.ToUpper(cMethod) == vMethod {

          temp = append(temp, vMethod)
          break
        }
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

  return "Content-Type" // fallback into default value
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
