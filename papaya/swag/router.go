package swag

import (
  "PapayaNet/papaya/koala/mapping"
)

type SwagRouter struct {
  composes []SwagComposeImpl
  group    *SwagGroup
}

type SwagRouterImpl interface {
  Init(group *SwagGroup)
  Group(name string) *SwagGroup
  Get(path string, expect *mapping.KMap, handler SwagHandler)
  Head(path string, expect *mapping.KMap, handler SwagHandler)
  Post(path string, expect *mapping.KMap, handler SwagHandler)
  Put(path string, expect *mapping.KMap, handler SwagHandler)
  Delete(path string, expect *mapping.KMap, handler SwagHandler)
  Connect(path string, expect *mapping.KMap, handler SwagHandler)
  Options(path string, expect *mapping.KMap, handler SwagHandler)
  Trace(path string, expect *mapping.KMap, handler SwagHandler)
}

func MakeSwagRouter(group *SwagGroup) *SwagRouter {

  router := &SwagRouter{}
  router.Init(group)

  return router
}

func (r *SwagRouter) Init(group *SwagGroup) {

  r.composes = make([]SwagComposeImpl, 0)
  r.group = group
}

func (r *SwagRouter) Group(name string) *SwagGroup {

  return MakeSwagGroup(r.group.Name + "/" + name)
}

func (r *SwagRouter) Get(path string, expect *mapping.KMap, handler SwagHandler) {

  compose := &SwagCompose{}
  compose.Init("GET", path, expect, handler)

  r.composes = append(r.composes, compose)
}

func (r *SwagRouter) Head(path string, expect *mapping.KMap, handler SwagHandler) {

  compose := &SwagCompose{}
  compose.Init("HEAD", path, expect, handler)

  r.composes = append(r.composes, compose)
}

func (r *SwagRouter) Post(path string, expect *mapping.KMap, handler SwagHandler) {

  compose := &SwagCompose{}
  compose.Init("POST", path, expect, handler)

  r.composes = append(r.composes, compose)
}

func (r *SwagRouter) Put(path string, expect *mapping.KMap, handler SwagHandler) {

  compose := &SwagCompose{}
  compose.Init("PUT", path, expect, handler)

  r.composes = append(r.composes, compose)
}

func (r *SwagRouter) Delete(path string, expect *mapping.KMap, handler SwagHandler) {

  compose := &SwagCompose{}
  compose.Init("DELETE", path, expect, handler)

  r.composes = append(r.composes, compose)
}

func (r *SwagRouter) Connect(path string, expect *mapping.KMap, handler SwagHandler) {

  compose := &SwagCompose{}
  compose.Init("CONNECT", path, expect, handler)

  r.composes = append(r.composes, compose)
}

func (r *SwagRouter) Options(path string, expect *mapping.KMap, handler SwagHandler) {

  compose := &SwagCompose{}
  compose.Init("OPTIONS", path, expect, handler)

  r.composes = append(r.composes, compose)
}

func (r *SwagRouter) Trace(path string, expect *mapping.KMap, handler SwagHandler) {

  compose := &SwagCompose{}
  compose.Init("TRACE", path, expect, handler)

  r.composes = append(r.composes, compose)
}
