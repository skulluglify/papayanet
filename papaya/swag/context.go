package swag

import "github.com/gofiber/fiber/v2"

type SwagContext struct {
  *fiber.Ctx
}

func MakeSwagContext(ctx *fiber.Ctx) *SwagContext {

  return &SwagContext{
    Ctx: ctx,
  }
}

func (c *SwagContext) Send(body []byte) error {

  return c.Ctx.Send(body)
}

func (c *SwagContext) SendFile(file string, compress ...bool) error {

  return c.Ctx.SendFile(file, compress...)
}
