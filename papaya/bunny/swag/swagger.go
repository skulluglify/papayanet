package swag

import (
  "net/http"
  "skfw/papaya/koala"
  "skfw/papaya/koala/collection"
  "skfw/papaya/koala/kornet"
  m "skfw/papaya/koala/mapping"
  "skfw/papaya/koala/tools/posix"
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
  SwagTasksQueueImpl
  version  koala.KVersionImpl
  info     *SwagInfo
  tag      string
  tags     []m.KMapImpl
  paths    m.KMapImpl
  root     posix.KPathImpl
  composes collection.KListImpl[SwagComposeImpl]
  renderer SwagRendererImpl
}

type SwagImpl interface {
  Init(app *fiber.App, info *SwagInfo)
  Version() koala.KVersionImpl
  Group(path string, tag string) SwagGroupImpl // alias as a Group('/')
  Router() SwagRouterImpl                      // alias as a Group('/')
  AddTag(tag string)
  AddPath(path string, method string, expect *SwagExpect)
  AddTask(task *SwagTask)
  Start() error
}

func MakeSwag(app *fiber.App, info *SwagInfo) SwagImpl {

  swag := &Swag{}
  swag.Init(app, info)

  return swag
}

func (swag *Swag) Init(app *fiber.App, info *SwagInfo) {

  swag.App = app
  swag.SwagTasksQueueImpl = SwagTasksQueueNew()

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

  swag.renderer = SwagRendererNew("/api/v3/openapi.json", app)
}

func (swag *Swag) Version() koala.KVersionImpl {

  return swag.version
}

func (swag *Swag) Group(path string, tag string) SwagGroupImpl {

  tag = swag.tag + "\\" + tag
  group := MakeSwagGroup(swag.root.Copy().Join(posix.KPathNew(path)), tag)
  group.Bind(swag.composes)

  return group
}

func (swag *Swag) Router() SwagRouterImpl {

  group := MakeSwagGroup(swag.root.Copy(), swag.tag)
  group.Bind(swag.composes)

  return group.Router()
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

  path = PathWrapSpecialNameWithBrackets(path)
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

func (swag *Swag) AddTask(task *SwagTask) {

  swag.SwagTasksQueueImpl.AddTask(task)
}

func (swag *Swag) Start() error {

  if err := swag.composes.ForEach(func(i uint, value SwagComposeImpl) error {

    method := value.Method()
    exp := value.Expect()
    tag := value.Tag()
    path := value.Path()
    expect := SwagExpectEvaluation(exp, []string{tag})

    requestValidation := expect.RequestValidation

    swag.Add(method, path, func(ctx *fiber.Ctx) error {

      // auth token
      // request validation

      context := MakeSwagContext(ctx, false)

      if requestValidation {

        validator := SwagRequestValidatorNew(exp, ctx)
        req, err := validator.Validation()

        if err != nil {

          return ctx.Status(http.StatusBadRequest).JSON(kornet.MessageNew(err.Error(), true))
        }

        context.Bind(req, nil)
      }

      if err := swag.SwagTasksQueueImpl.Start(exp, context); err != nil {

        if !context.Revoke() {

          return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
        }

        return nil
      }

      return value.Handler(context)
    })

    if !m.KValueToBool(exp.Get("hidden")) {

      swag.AddTag(tag)
      swag.AddPath(path, method, expect)
    }

    return nil

  }); err != nil {

    return err
  }

  return swag.renderer.Render(swag.version, swag.info, swag.tags, swag.paths)
}
