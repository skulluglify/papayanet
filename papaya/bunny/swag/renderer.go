package swag

import (
  "PapayaNet/papaya/bunny/swag/binaries"
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

func SwagTemplateHandler(ctx *fiber.Ctx) error {

  req := ctx.Request()
  URI := req.URI()
  url := &url.URL{
    Scheme: string(URI.Scheme()),
    Host:   string(URI.Host()),
    Path:   "/api/v3/openapi.json",
  }

  swagTemplate := []byte(SwagTemplateHTML(url.String()))

  ctx.Set("Content-Type", "text/html")
  ctx.Set("Content-Length", strconv.Itoa(len(swagTemplate)))

  return ctx.Send(swagTemplate)
}

func SwagRedocTemplateHandler(ctx *fiber.Ctx) error {

  req := ctx.Request()
  URI := req.URI()
  url := &url.URL{
    Scheme: string(URI.Scheme()),
    Host:   string(URI.Host()),
    Path:   "/api/v3/openapi.json",
  }

  swagTemplate := []byte(SwagRedocTemplateHTML(url.String()))

  ctx.Set("Content-Type", "text/html")
  ctx.Set("Content-Length", strconv.Itoa(len(swagTemplate)))

  return ctx.Send(swagTemplate)
}
