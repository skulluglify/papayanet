package swag

import (
  "net/http"
  "skfw/papaya/bunny/swag/method"
  "skfw/papaya/koala/collection"
  "skfw/papaya/koala/mapping"
  "skfw/papaya/koala/tools/posix"
)

type SwagRouter struct {
  composes collection.KListImpl[SwagComposeImpl]
  path     posix.KPathImpl
  tag      string
}

type SwagRouterImpl interface {
  Init(group SwagGroupImpl)
  Group(path string, tag string) SwagGroupImpl
  Get(path string, expect *mapping.KMap, handler SwagRouteHandler)
  Head(path string, expect *mapping.KMap, handler SwagRouteHandler)
  Post(path string, expect *mapping.KMap, handler SwagRouteHandler)
  Put(path string, expect *mapping.KMap, handler SwagRouteHandler)
  Delete(path string, expect *mapping.KMap, handler SwagRouteHandler)
  Connect(path string, expect *mapping.KMap, handler SwagRouteHandler)
  Options(path string, expect *mapping.KMap, handler SwagRouteHandler)
  Trace(path string, expect *mapping.KMap, handler SwagRouteHandler)
  Bind(composes collection.KListImpl[SwagComposeImpl])
  Composes() collection.KListImpl[SwagComposeImpl]
}

func MakeSwagRouter(group *SwagGroup) SwagRouterImpl {

  router := &SwagRouter{}
  router.Init(group)

  return router
}

func (router *SwagRouter) Init(group SwagGroupImpl) {

  router.composes = MakeSwagComposes()
  router.path = group.Path()
  router.tag = group.Tag()
}

func (router *SwagRouter) Group(path string, tag string) SwagGroupImpl {

  tag = router.tag + "\\" + tag
  group := MakeSwagGroup(router.path.Join(posix.KPathNew(path)), tag)
  group.Bind(router.composes)

  return group
}

func (router *SwagRouter) OptionsBypass(path string) {

  router.Options(path, &mapping.KMap{
    "hidden":      true,
    "description": "Passing Options On Chrome Browser",
  }, func(ctx *SwagContext) error {

    return ctx.Status(http.StatusOK).Send(nil)
  })
}

func (router *SwagRouter) Get(path string, expect *mapping.KMap, handler SwagRouteHandler) {

  p := router.path.Copy()
  if path != "" {
    p = p.Join(posix.KPathNew(path))
  }

  compose := MakeSwagCompose(method.GET, router.tag, p, expect, handler)
  router.composes.Push(compose)

  // fix issue options
  router.OptionsBypass(path)
}

func (router *SwagRouter) Head(path string, expect *mapping.KMap, handler SwagRouteHandler) {

  p := router.path.Copy()
  if path != "" {
    p = p.Join(posix.KPathNew(path))
  }

  compose := MakeSwagCompose(method.HEAD, router.tag, p, expect, handler)
  router.composes.Push(compose)

  // fix issue options
  router.OptionsBypass(path)
}

func (router *SwagRouter) Post(path string, expect *mapping.KMap, handler SwagRouteHandler) {

  p := router.path.Copy()
  if path != "" {
    p = p.Join(posix.KPathNew(path))
  }

  compose := MakeSwagCompose(method.POST, router.tag, p, expect, handler)
  router.composes.Push(compose)

  // fix issue options
  router.OptionsBypass(path)
}

func (router *SwagRouter) Put(path string, expect *mapping.KMap, handler SwagRouteHandler) {

  p := router.path.Copy()
  if path != "" {
    p = p.Join(posix.KPathNew(path))
  }

  compose := MakeSwagCompose(method.PUT, router.tag, p, expect, handler)
  router.composes.Push(compose)

  // fix issue options
  router.OptionsBypass(path)
}

func (router *SwagRouter) Delete(path string, expect *mapping.KMap, handler SwagRouteHandler) {

  p := router.path.Copy()
  if path != "" {
    p = p.Join(posix.KPathNew(path))
  }

  compose := MakeSwagCompose(method.DELETE, router.tag, p, expect, handler)
  router.composes.Push(compose)

  // fix issue options
  router.OptionsBypass(path)
}

func (router *SwagRouter) Connect(path string, expect *mapping.KMap, handler SwagRouteHandler) {

  p := router.path.Copy()
  if path != "" {
    p = p.Join(posix.KPathNew(path))
  }

  compose := MakeSwagCompose(method.CONNECT, router.tag, p, expect, handler)
  router.composes.Push(compose)

  // fix issue options
  router.OptionsBypass(path)
}

func (router *SwagRouter) Options(path string, expect *mapping.KMap, handler SwagRouteHandler) {

  p := router.path.Copy()
  if path != "" {
    p = p.Join(posix.KPathNew(path))
  }

  compose := MakeSwagCompose(method.OPTIONS, router.tag, p, expect, handler)
  router.composes.Push(compose)
}

func (router *SwagRouter) Trace(path string, expect *mapping.KMap, handler SwagRouteHandler) {

  p := router.path.Copy()
  if path != "" {
    p = p.Join(posix.KPathNew(path))
  }

  compose := MakeSwagCompose(method.TRACE, router.tag, p, expect, handler)
  router.composes.Push(compose)

  // fix issue options
  router.OptionsBypass(path)
}

func (router *SwagRouter) Bind(composes collection.KListImpl[SwagComposeImpl]) {

  router.composes = composes
}

func (router *SwagRouter) Composes() collection.KListImpl[SwagComposeImpl] {

  return router.composes
}
