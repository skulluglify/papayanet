package swag

import (
	"PapayaNet/papaya/koala/collection"
	"PapayaNet/papaya/koala/tools/posix"
)

type SwagGroup struct {
	composes collection.KListImpl[SwagComposeImpl]
	tag      string
	path     posix.KPathImpl
}

type SwagGroupImpl interface {
	Init(path posix.KPathImpl, tag string)
	Group(path string, tag string) SwagGroupImpl
	Router() SwagRouterImpl
	BindingComposes(composes collection.KListImpl[SwagComposeImpl])
	Composes() collection.KListImpl[SwagComposeImpl]
	Path() posix.KPathImpl
	Tag() string
}

func MakeSwagGroup(path posix.KPathImpl, tag string) SwagGroupImpl {

	group := &SwagGroup{}
	group.Init(path, tag)

	return group
}

func (group *SwagGroup) Init(path posix.KPathImpl, tag string) {

	group.composes = MakeSwagComposes()
	group.tag = tag
	group.path = path
}

func (group *SwagGroup) Group(path string, tag string) SwagGroupImpl {

	tag = group.tag + "\\" + tag
	swagGroup := MakeSwagGroup(group.path.Join(posix.KPathNew(path)), tag)
	swagGroup.BindingComposes(group.composes)

	return swagGroup
}

func (group *SwagGroup) Router() SwagRouterImpl {

	router := MakeSwagRouter(group)
	router.BindingComposes(group.composes)

	return router
}

func (group *SwagGroup) BindingComposes(composes collection.KListImpl[SwagComposeImpl]) {

	group.composes = composes
}

func (group *SwagGroup) Composes() collection.KListImpl[SwagComposeImpl] {

	return group.composes
}

func (group *SwagGroup) Path() posix.KPathImpl {

	return group.path
}

func (group *SwagGroup) Tag() string {

	return group.tag
}
