package swag

import (
  "PapayaNet/papaya/koala/collection"
  "PapayaNet/papaya/koala/mapping"
  "PapayaNet/papaya/koala/tools/posix"
)

// compose, from expectation

type SwagCompose struct {
  method  string
  tag     string
  path    posix.KPathImpl
  expect  mapping.KMapImpl
  handler SwagRouteHandler
}

type SwagComposeImpl interface {
  Init(method string, tag string, path posix.KPathImpl, expect mapping.KMapImpl, handler SwagRouteHandler)
  Method() string
  Tag() string
  Path() string
  Expect() mapping.KMapImpl
  Handler(ctx *SwagContext) error
}

func MakeSwagComposes() collection.KListImpl[SwagComposeImpl] {

  return collection.KListNew[SwagComposeImpl]()
}

func MakeSwagCompose(method string, tag string, path posix.KPathImpl, expect mapping.KMapImpl, handler SwagRouteHandler) SwagComposeImpl {

  compose := &SwagCompose{}
  compose.Init(method, tag, path, expect, handler)

  return compose
}

func (c *SwagCompose) Init(method string, tag string, path posix.KPathImpl, expect mapping.KMapImpl, handler SwagRouteHandler) {

  c.method = method
  c.tag = tag
  c.path = path
  c.expect = expect
  c.handler = handler
}

func (c *SwagCompose) Method() string {

  return c.method
}

func (c *SwagCompose) Tag() string {

  return c.tag
}

func (c *SwagCompose) Path() string {

  // normalize path, make it compatible in openapi 3
  return c.path.String()
}

func (c *SwagCompose) Expect() mapping.KMapImpl {

  return c.expect
}

func (c *SwagCompose) Handler(ctx *SwagContext) error {

  return c.handler(ctx)
}
