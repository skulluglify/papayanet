package swag

import (
  "PapayaNet/papaya/koala"
  "PapayaNet/papaya/koala/mapping"
  "github.com/gofiber/fiber/v2"
)

const (
  VersionMajor = 2
  VersionMinor = 0
  VersionPatch = 0
)

type Swag struct {
  *fiber.App
  version koala.KVersionImpl
  info    *SwagInfo
  path    mapping.KMap
  group   SwagGroupImpl
}

type SwagImpl interface {
  Init(app *fiber.App, info *SwagInfo)
  Version() koala.KVersionImpl
  Router() SwagRouterImpl // alias as a Group('/')
  Swagger() mapping.KMap
  Start()
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
}

func (swag *Swag) Version() koala.KVersionImpl {

  return swag.version
}

func (swag *Swag) Router() SwagRouterImpl {

  swag.group = MakeSwagGroup("/")

  return swag.group.Router()
}

func (swag *Swag) Swagger() mapping.KMap {
  //TODO implement me
  panic("implement me")
}

func (swag *Swag) Start() {
  //TODO implement me
  panic("implement me")
}
