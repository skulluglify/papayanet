package swag

import (
  "PapayaNet/papaya/bunny/swag/binaries"
  "PapayaNet/papaya/koala"
  "PapayaNet/papaya/koala/kio/leaf"
  m "PapayaNet/papaya/koala/mapping"
  "PapayaNet/papaya/koala/pp"
  "fmt"
  "net/url"
  "strconv"

  "github.com/gofiber/fiber/v2"
)

type SwagFileAssets struct {
  SwagUIStyle      string `url:"https://unpkg.com/swagger-ui-dist@latest/swagger-ui.css"`
  SwagBundleScript string `url:"https://unpkg.com/swagger-ui-dist@latest/swagger-ui-bundle.js"`
  SwagRedocScript  string `url:"https://rebilly.github.io/ReDoc/releases/latest/redoc.min.js"`
}

// func SwagTemplateHTML(url string) string {

//   return fmt.Sprintf(`
//   <!DOCTYPE html>
//   <html lang="en">
//     <head>
//       <meta charset="utf-8" />
//       <meta name="viewport" content="width=device-width, initial-scale=1" />
//       <meta
//         name="description"
//         content="SwaggerUI"
//       />
//       <title>SwaggerUI</title>
//       <link rel="stylesheet" href="%s" />
//     </head>
//     <body>
//     <div id="swagger-ui"></div>
//     <script src="%s" crossorigin></script>
//     <script src="%s" crossorigin></script>
//     <script>
//       window.onload = () => {
//         window.ui = SwaggerUIBundle({
//           url: '%s',
//           dom_id: '#swagger-ui',
//           presets: [
//             SwaggerUIBundle.presets.apis,
//             SwaggerUIStandalonePreset
//           ],
//           layout: "StandaloneLayout",
//         });
//       };
//     </script>
//     </body>
//   </html>
//   `, binaries.SwagBinStyles, binaries.SwagBinScripts, binaries.SwagBinPresets, url)
// }

func SwagTemplateHTML(url string) string {

  return fmt.Sprintf(`
  <!DOCTYPE html>
  <html lang="en">
    <head>
      <meta charset="utf-8" />
      <meta name="viewport" content="width=device-width, initial-scale=1" />
      <meta
        name="description"
        content="SwaggerUI"
      />
      <title>SwaggerUI</title>
      <link rel="stylesheet" href="%s" />
    </head>
    <body>
    <div id="swagger-ui"></div>
    <script src="%s" crossorigin></script>
    <script>
      window.onload = () => {
        window.ui = SwaggerUIBundle({
          url: '%s',
          dom_id: '#swagger-ui',
        });
      };
    </script>
    </body>
  </html>
  `, binaries.SwagBinStyles, binaries.SwagBinScripts, url)
}

func SwagRedocTemplateHTML(url string) string {

  return fmt.Sprintf(`
  <!DOCTYPE html>
  <html lang="en">
    <head>
      <meta charset="utf-8" />
      <meta name="viewport" content="width=device-width, initial-scale=1" />
      <meta
        name="description"
        content="Redoc"
      />
      <title>Redoc</title>
    </head>
    <body>
    <redoc spec-url="%s"></redoc>
    <script src="%s"></script>
    </body>
  </html>
  `, url, binaries.SwagBinRedocScripts)
}

type SwagRenderer struct {
  *fiber.App
  SwagUIHTML leaf.KBufferImpl // cache
  RedocHTML  leaf.KBufferImpl
  OpenAPI    leaf.KBufferImpl
  Path       string
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
  r.SwagUIHTML = leaf.KMakeBufferZone(0)
  r.RedocHTML = leaf.KMakeBufferZone(0)
  r.OpenAPI = leaf.KMakeBufferZone(0)

  r.Path = path
  r.App = app
}

func (r *SwagRenderer) GetRequestOpenApi(ctx *fiber.Ctx) string {

  req := ctx.Request()
  URI := req.URI()
  u := &url.URL{
    Scheme: string(URI.Scheme()),
    Host:   string(URI.Host()),
    Path:   pp.QStr(r.Path, "/api/v3/openapi.json"),
  }

  return u.String()
}

func (r *SwagRenderer) Render(version koala.KVersionImpl, info *SwagInfo, tags []m.KMapImpl, paths m.KMapImpl) error {

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

  r.Get(pp.QStr(r.Path, "/api/v3/openapi.json"), func(ctx *fiber.Ctx) error {

    r.OpenAPI.Seek(0)

    ctx.Set("Content-Type", "application/json")
    ctx.Set("Content-Length", strconv.FormatUint(uint64(r.OpenAPI.Size()), 10))

    return ctx.Send(r.OpenAPI.ReadAll())
  })

  r.Get("/swag", func(ctx *fiber.Ctx) error {

    // reset seek point
    r.SwagUIHTML.Seek(0)

    // lazy load
    if r.SwagUIHTML.Size() == 0 {

      openAPI := r.GetRequestOpenApi(ctx)
      r.SwagUIHTML = leaf.KMakeBuffer([]byte(SwagTemplateHTML(openAPI)))
    }

    ctx.Set("Content-Type", "text/html")
    ctx.Set("Content-Length", strconv.FormatUint(uint64(r.SwagUIHTML.Size()), 10))

    return ctx.Send(r.SwagUIHTML.ReadAll())
  })

  r.Get("/redoc", func(ctx *fiber.Ctx) error {

    // reset seek point
    r.RedocHTML.Seek(0)

    // lazy load
    if r.RedocHTML.Size() == 0 {

      openAPI := r.GetRequestOpenApi(ctx)
      r.RedocHTML = leaf.KMakeBuffer([]byte(SwagRedocTemplateHTML(openAPI)))
    }

    ctx.Set("Content-Type", "text/html")
    ctx.Set("Content-Length", strconv.FormatUint(uint64(r.RedocHTML.Size()), 10))

    return ctx.Send(r.RedocHTML.ReadAll())
  })

  return nil
}

func (r *SwagRenderer) Close() error {

  r.SwagUIHTML.Close()
  r.RedocHTML.Close()

  return nil
}
