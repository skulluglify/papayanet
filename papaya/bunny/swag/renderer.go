package swag

import (
  "net/url"
  "skfw/papaya/ant/bpack"
  "skfw/papaya/koala"
  "skfw/papaya/koala/kio/leaf"
  m "skfw/papaya/koala/mapping"
  "skfw/papaya/koala/pp"
  "strconv"

  "github.com/gofiber/fiber/v2"
)

type SwagRenderer struct {
  *fiber.App
  SwagPkt        *bpack.Packet
  RedocPkt       *bpack.Packet
  SwagStylePkt   *bpack.Packet
  SwagScriptPkt  *bpack.Packet
  SwagPresetPkt  *bpack.Packet
  RedocScriptPkt *bpack.Packet
  OpenAPI        leaf.KBufferImpl
  Path           string
}

type SwagRendererImpl interface {
  Init(path string, app *fiber.App)
  GetRequestOpenApi(ctx *fiber.Ctx) string
  Render(version koala.KVersionImpl, info *SwagInfo, tags []m.KMapImpl, paths m.KMapImpl) error
  Close() error
}

func SwagRendererNew(path string, app *fiber.App) SwagRendererImpl {

  renderer := &SwagRenderer{}
  renderer.Init(path, app)
  return renderer
}

func (r *SwagRenderer) Init(path string, app *fiber.App) {

  // set zero value
  r.OpenAPI = leaf.KMakeBufferZone(0)

  r.SwagPkt = bpack.OpenPacket("/data/swag/index.html")
  r.RedocPkt = bpack.OpenPacket("/data/swag/redoc.html")
  r.SwagStylePkt = bpack.OpenPacket("/data/swag/ui.css")
  r.SwagScriptPkt = bpack.OpenPacket("/data/swag/swagger.js")
  r.SwagPresetPkt = bpack.OpenPacket("/data/swag/preset.js")
  r.RedocScriptPkt = bpack.OpenPacket("/data/swag/redoc.js")

  r.Path = path
  r.App = app
}

func (r *SwagRenderer) GetRequestOpenApi(ctx *fiber.Ctx) string {

  req := ctx.Request()
  URI := req.URI()
  u := &url.URL{
    Scheme: string(URI.Scheme()),
    Host:   string(URI.Host()),
    Path:   pp.Qstr(r.Path, "/api/v3/openapi.json"),
  }

  return u.String()
}

func (r *SwagRenderer) Render(version koala.KVersionImpl, info *SwagInfo, tags []m.KMapImpl, paths m.KMapImpl) error {

  var err error

  data := &m.KMap{
    "openapi": version.String(),
    "info": &m.KMap{
      "title":       info.Title,
      "description": info.Description,
      "version":     info.Version,
    },
    "tags":  tags,
    "paths": paths,
  }

  // cache on temporary
  r.OpenAPI = leaf.KMakeBuffer([]byte(data.JSON()))

  r.App.Get(pp.Qstr(r.Path, "/api/v3/openapi.json"), func(ctx *fiber.Ctx) error {

    r.OpenAPI.Seek(0)

    ctx.Set("Content-Type", "application/json")
    ctx.Set("Content-Length", strconv.FormatUint(uint64(r.OpenAPI.Size()), 10))

    return ctx.Send(r.OpenAPI.ReadAll())
  })

  if err = bpack.HttpExposePacket(r.App, "/doc/swag", "text/html", r.SwagPkt); err != nil {
    return err
  }

  if err = bpack.HttpExposePacket(r.App, "/doc/swag/redoc", "text/html", r.RedocPkt); err != nil {
    return err
  }

  return nil
}

func (r *SwagRenderer) Close() error {

  return nil
}
