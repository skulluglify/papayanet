package swag

import (
  "PapayaNet/papaya/koala"
  "PapayaNet/papaya/koala/collection"
  m "PapayaNet/papaya/koala/mapping"
  "PapayaNet/papaya/koala/tools/posix"
  "strconv"
  "strings"

  "github.com/gofiber/fiber/v2"
)

// using openapi 3
const (
  VersionMajor = 3
  VersionMinor = 0
  VersionPatch = 0
)

type Swag struct {
  *fiber.App
  version  koala.KVersionImpl
  info     *SwagInfo
  tag      string
  tags     []m.KMapImpl
  paths    m.KMapImpl
  root     posix.KPathImpl
  composes collection.KListImpl[SwagComposeImpl] // can be push item without passing re-value
}

type SwagImpl interface {
  Init(app *fiber.App, info *SwagInfo)
  Version() koala.KVersionImpl
  Group(path string, tag string) SwagGroupImpl // alias as a Group('/')
  Router() SwagRouterImpl                      // alias as a Group('/')
  AddTag(tag string)
  AddPath(path string, method string, expect *SwagExpect)
  Start() error
  Swagger()
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

  swag.tags = make([]m.KMapImpl, 0)
  swag.paths = &m.KMap{}

  swag.root = posix.KPathNew("/")
  swag.composes = MakeSwagComposes()
}

func (swag *Swag) Version() koala.KVersionImpl {

  return swag.version
}

func (swag *Swag) Group(path string, tag string) SwagGroupImpl {

  tag = swag.tag + "\\" + tag
  group := MakeSwagGroup(swag.root.Join(posix.KPathNew(path)), tag)
  group.Bind(swag.composes)

  return group
}

func (swag *Swag) Router() SwagRouterImpl {

  group := MakeSwagGroup(swag.root, swag.tag)
  group.Bind(swag.composes)

  return group.Router()
}

func (swag *Swag) Swagger() {

  data := &m.KMap{
    "openapi": swag.version.String(),
    "info": &m.KMap{
      "title":       swag.info.Title,
      "description": swag.info.Description,
      "version":     swag.info.Version,
    },
    "tags":  swag.tags,
    "paths": swag.paths,
  }

  // cache on temporary
  temp := []byte(data.JSON())

  swag.Get("/api/v3/openapi.json", func(ctx *fiber.Ctx) error {

    ctx.Set("Content-Type", "application/json")
    ctx.Set("Content-Length", strconv.Itoa(len(temp)))

    return ctx.Send(temp)
  })

  swag.Get("/swag", SwagTemplateHandler)
  swag.Get("/redoc", SwagRedocTemplateHandler)
}

func (swag *Swag) AddTag(tag string) {

  dupTag := false

  for _, dTag := range swag.tags {

    if t := m.KValueToString(dTag.Get("name")); t != "" {

      if t == tag {

        dupTag = true
        break
      }
    }
  }

  if !dupTag {

    swag.tags = append(swag.tags, &m.KMap{
      "name": tag,
    })
  }
}

func (swag *Swag) AddPath(path string, method string, expect *SwagExpect) {

  var data m.KMapImpl

  path = SwagPathNorm(path)
  data = nil

  if currPath := swag.paths.Get(path); currPath != nil {

    if mm := m.KMapCast(currPath); mm != nil {

      data = mm
    }
  }

  if data != nil {

    data.Put(strings.ToLower(method), expect.Path)

  } else {

    data = &m.KMap{}
    data.Put(strings.ToLower(method), expect.Path)
    swag.paths.Put(path, data)
  }
}

func (swag *Swag) Start() error {

  if err := swag.composes.ForEach(func(i uint, value SwagComposeImpl) error {

    method := value.Method()
    tag := value.Tag()
    path := value.Path()
    expect := SwagExpectEvaluation(value.Expect(), []string{tag})

    authToken := expect.AuthToken
    requestValidation := expect.RequestValidation

    println(authToken, requestValidation)

    swag.Add(method, path, func(ctx *fiber.Ctx) error {

      // auth token
      // request validation

      return value.Handler(&SwagContext{
        Ctx: ctx,
      })
    })

    swag.AddTag(tag)
    swag.AddPath(path, method, expect)

    return nil

  }); err != nil {

    return err
  }

  return nil
}
