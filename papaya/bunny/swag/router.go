package swag

import (
  "PapayaNet/papaya/bunny/swag/method"
  "PapayaNet/papaya/koala/collection"
  "PapayaNet/papaya/koala/mapping"
  "PapayaNet/papaya/koala/tools/posix"
)

type SwagRouter struct {
  composes collection.KListImpl[SwagComposeImpl]
  path     posix.KPathImpl
  tag      string
}

type SwagRouterImpl interface {
  Init(group SwagGroupImpl)
  Group(path string, tag string) SwagGroupImpl
  Get(path string, expect *mapping.KMap, handler SwagHandler)
  Head(path string, expect *mapping.KMap, handler SwagHandler)
  Post(path string, expect *mapping.KMap, handler SwagHandler)
  Put(path string, expect *mapping.KMap, handler SwagHandler)
  Delete(path string, expect *mapping.KMap, handler SwagHandler)
  Connect(path string, expect *mapping.KMap, handler SwagHandler)
  Options(path string, expect *mapping.KMap, handler SwagHandler)
  Trace(path string, expect *mapping.KMap, handler SwagHandler)
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

func (router *SwagRouter) Get(path string, expect *mapping.KMap, handler SwagHandler) {

  p := router.path.Copy()
  if path != "" {
    p = p.Join(posix.KPathNew(path))
  }

  compose := MakeSwagCompose(method.GET, router.tag, p, expect, handler)
  router.composes.Push(compose)
}

func (router *SwagRouter) Head(path string, expect *mapping.KMap, handler SwagHandler) {

  p := router.path.Copy()
  if path != "" {
    p = p.Join(posix.KPathNew(path))
  }

  compose := MakeSwagCompose(method.HEAD, router.tag, p, expect, handler)
  router.composes.Push(compose)
}

func (router *SwagRouter) Post(path string, expect *mapping.KMap, handler SwagHandler) {

  p := router.path.Copy()
  if path != "" {
    p = p.Join(posix.KPathNew(path))
  }

  compose := MakeSwagCompose(method.POST, router.tag, p, expect, handler)
  router.composes.Push(compose)
}

func (router *SwagRouter) Put(path string, expect *mapping.KMap, handler SwagHandler) {

  p := router.path.Copy()
  if path != "" {
    p = p.Join(posix.KPathNew(path))
  }

  compose := MakeSwagCompose(method.PUT, router.tag, p, expect, handler)
  router.composes.Push(compose)
}

func (router *SwagRouter) Delete(path string, expect *mapping.KMap, handler SwagHandler) {

  p := router.path.Copy()
  if path != "" {
    p = p.Join(posix.KPathNew(path))
  }

  compose := MakeSwagCompose(method.DELETE, router.tag, p, expect, handler)
  router.composes.Push(compose)
}

func (router *SwagRouter) Connect(path string, expect *mapping.KMap, handler SwagHandler) {

  p := router.path.Copy()
  if path != "" {
    p = p.Join(posix.KPathNew(path))
  }

  compose := MakeSwagCompose(method.CONNECT, router.tag, p, expect, handler)
  router.composes.Push(compose)
}

func (router *SwagRouter) Options(path string, expect *mapping.KMap, handler SwagHandler) {

  p := router.path.Copy()
  if path != "" {
    p = p.Join(posix.KPathNew(path))
  }

  compose := MakeSwagCompose(method.OPTIONS, router.tag, p, expect, handler)
  router.composes.Push(compose)
}

func (router *SwagRouter) Trace(path string, expect *mapping.KMap, handler SwagHandler) {

  p := router.path.Copy()
  if path != "" {
    p = p.Join(posix.KPathNew(path))
  }

  compose := MakeSwagCompose(method.TRACE, router.tag, p, expect, handler)
  router.composes.Push(compose)
}

func (router *SwagRouter) Bind(composes collection.KListImpl[SwagComposeImpl]) {

  router.composes = composes
}

func (router *SwagRouter) Composes() collection.KListImpl[SwagComposeImpl] {

  return router.composes
}
