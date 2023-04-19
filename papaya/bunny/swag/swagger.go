package swag

import (
	"PapayaNet/papaya/koala"
	"PapayaNet/papaya/koala/collection"
	m "PapayaNet/papaya/koala/mapping"
	"PapayaNet/papaya/koala/tools/posix"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

const (
	VersionMajor = 2
	VersionMinor = 0
	VersionPatch = 0
)

type Swag struct {
	*fiber.App
	version  koala.KVersionImpl
	info     *SwagInfo
	tag      string
	paths    m.KMapImpl
	tags     []m.KMapImpl
	root     posix.KPathImpl
	composes collection.KListImpl[SwagComposeImpl]
}

type SwagImpl interface {
	Init(app *fiber.App, info *SwagInfo)
	Version() koala.KVersionImpl
	Group(path string, tag string) SwagGroupImpl // alias as a Group('/')
	Router() SwagRouterImpl                      // alias as a Group('/')
	Swagger() m.KMapImpl
	Start() error
}

func MakeSwag(app *fiber.App, info *SwagInfo) SwagImpl {

	swag := &Swag{}
	swag.Init(app, info)

	return swag
}

func (swag *Swag) Init(app *fiber.App, info *SwagInfo) {

	swag.App = app
	swag.info = info
	swag.version = koala.KVersionNew(
		VersionMajor,
		VersionMinor,
		VersionPatch,
	)
	swag.tag = "App"
	swag.root = posix.KPathNew("/")
	swag.composes = MakeSwagComposes()
}

func (swag *Swag) Version() koala.KVersionImpl {

	return swag.version
}

func (swag *Swag) Group(path string, tag string) SwagGroupImpl {

	tag = swag.tag + "\\" + tag
	group := MakeSwagGroup(swag.root.Join(posix.KPathNew(path)), tag)
	group.BindingComposes(swag.composes)

	return group
}

func (swag *Swag) Router() SwagRouterImpl {

	group := MakeSwagGroup(swag.root, swag.tag)
	group.BindingComposes(swag.composes)

	return group.Router()
}

func (swag *Swag) Swagger() m.KMapImpl {

	data := &m.KMap{
		"openapi": "3.0.0",
		"info": &m.KMap{
			"title":       swag.info.Title,
			"description": swag.info.Description,
			"version":     swag.info.Version,
		},
		"tags":  swag.tags,
		"paths": swag.paths,
	}

	return data
}

func (swag *Swag) Start() error {

	if err := swag.composes.ForEach(func(i uint, value SwagComposeImpl) error {

		method := value.Method()
		tag := value.Tag()
		path := value.Path()
		expect := value.Expect()

		fmt.Println(SwagContentFormatter(expect.Get("request.application/json.schema")).JSON())

		fmt.Println(i, method)
		fmt.Println(i, tag)
		fmt.Println(i, path)
		fmt.Println(i, expect)

		return nil
	}); err != nil {

		return err
	}

	return nil
}
