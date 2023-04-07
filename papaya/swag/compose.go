package swag

import (
  "PapayaNet/papaya/koala/mapping"
)

// compose, from expectation

type SwagCompose struct {
  method  string
  path    string
  expect  *mapping.KMap
  handler SwagHandler
}

type SwagComposeImpl interface {
  Init(method string, path string, expect *mapping.KMap, handler SwagHandler)
  Method() string
  Path() string
  Expect() *mapping.KMap
  Handler() SwagHandler
}

func (c *SwagCompose) Init(method string, path string, expect *mapping.KMap, handler SwagHandler) {

  c.method = method
  c.path = path
  c.expect = expect
  c.handler = handler
}

func (c *SwagCompose) Method() string {

  return c.method
}

func (c *SwagCompose) Path() string {

  return c.path
}

func (c *SwagCompose) Expect() *mapping.KMap {

  return c.expect
}

func (c *SwagCompose) Handler() SwagHandler {

  return c.handler
}
