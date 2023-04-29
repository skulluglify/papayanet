package cors

import m "skfw/papaya/bunny/swag/method"

var Methods = []string{
  m.GET,
  m.HEAD,
  m.POST,
  m.PUT,
  m.DELETE,
  m.CONNECT,
  m.OPTIONS,
  m.TRACE,
}

var Headers = []string{
  "Access-Control-Allow-Origin",
  "Access-Control-Allow-Methods",
  "Access-Control-Allow-Headers",
  "Access-Control-Allow-Credentials",
  "Access-Control-Max-Age",
}
