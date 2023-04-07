package swag

type SwagGroup struct {
  Name   string
  router *SwagRouter
}

type SwagGroupImpl interface {
  Init(name string)
  Router() *SwagRouter
}

func MakeSwagGroup(name string) *SwagGroup {

  group := &SwagGroup{}
  group.Init(name)

  return group
}

func (g *SwagGroup) Init(name string) {

  g.Name = name

  router := &SwagRouter{}
  router.Init(g)

  g.router = router
}

func (g *SwagGroup) Router() *SwagRouter {

  return g.router
}
