package swag

import "github.com/gofiber/fiber/v2"

type SwagContext struct {
  *fiber.Ctx
  event  bool
  data   any
  Closed bool
}

func MakeSwagContext(ctx *fiber.Ctx, event bool) *SwagContext {

  return &SwagContext{
    Ctx:   ctx,
    event: event,
    data:  nil,
  }
}

func MakeSwagContextWithEvent(ctx *fiber.Ctx, data any) *SwagContext {

  return &SwagContext{
    Ctx:   ctx,
    event: true,
    data:  data,
  }
}

func (c *SwagContext) Event() any {

  if c.event {

    if c.data == nil {

      return true
    }

    return c.data
  }

  return nil
}
