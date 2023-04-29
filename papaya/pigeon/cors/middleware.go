package cors

import (
  "net/http"

  "github.com/gofiber/fiber/v2"
)

func MakeMiddlewareForManageConsumers(manageConsumers ManageConsumersImpl) fiber.Handler {

  return func(ctx *fiber.Ctx) error {

    req, res := ctx.Request(), ctx.Response()

    // GET/HEAD not have body

    method := ctx.Method() // get method from request
    origin := string(req.Header.Peek("Origin"))

    // TODO: for features on future
    // AccessControlRequestMethod := string(req.Header.Peek("Access-Control-Request-Method")) // GET, POST, ...
    // AccessControlRequestPrivateNetwork := string(req.Header.Peek("Access-Control-Request-Private-Network")) // true or false
    // AccessControlRequestHeaders := string(req.Header.Peek("Access-Control-Request-Headers")) // Content-Type, ...

    var consumer ConsumerImpl

    // only for consumer
    if origin != "" {

      if consumer = manageConsumers.Get(method, origin); consumer != nil {

        // fixing a problem if origin a have prefix of www.
        if consumer = consumer.Origin(origin); consumer != nil {

          consumer.Apply(res)

          // consumer
          return ctx.Next()
        }
      }

      // blocked
      return ctx.Status(http.StatusNoContent).SendString("") // stop processing
    }

    // self-hosted
    return ctx.Next()
  }
}
