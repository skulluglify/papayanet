package swag

import (
  "github.com/gofiber/fiber/v2"
  "skfw/papaya/koala/kornet"
)

type SwagContext struct {
  *fiber.Ctx
  event    bool
  data     any
  body     []byte
  kReq     *kornet.Request
  kRes     *kornet.Response
  finished bool
}

type SwagContextImpl interface {
  Event() bool    // if a task running, event is true
  Target() any    // Event
  Solve(data any) // --> Target()
  Modify(body []byte)
  Body() []byte
  Status(status int) *SwagContext
  Prevent()
  Revoke() bool
  Kornet() (*kornet.Request, *kornet.Response)
  Bind(req *kornet.Request, res *kornet.Response)
}

func MakeSwagContext(ctx *fiber.Ctx, event bool) *SwagContext {

  return &SwagContext{
    Ctx:   ctx,
    event: event,
    data:  nil,
  }
}

func MakeSwagContextEvent(ctx *fiber.Ctx, data any) *SwagContext {

  return &SwagContext{
    Ctx:   ctx,
    event: true,
    body:  nil,
    data:  data,
  }
}

func (c *SwagContext) Event() bool {

  return c.event
}

func (c *SwagContext) Target() any {

  return c.data
}

func (c *SwagContext) Solve(data any) {

  c.event = true // solve value from key expectation
  c.data = data  // value from expectation
}

func (c *SwagContext) Modify(body []byte) {

  c.body = body
}

// wrapping Body with Modify Body

func (c *SwagContext) Body() []byte {

  if c.body != nil {

    return c.body
  }

  return c.Ctx.Body()
}

// wrapping Status with New SwagContext

func (c *SwagContext) Status(status int) *SwagContext {

  return &SwagContext{
    Ctx:      c.Ctx.Status(status),
    event:    c.event,
    data:     c.data,
    body:     c.body,
    finished: c.finished,
  }
}

func (c *SwagContext) Prevent() {

  c.finished = true
}

func (c *SwagContext) Revoke() bool {

  return c.finished
}

func (c *SwagContext) Kornet() (*kornet.Request, *kornet.Response) {

  return c.kReq, c.kRes
}

func (c *SwagContext) Bind(req *kornet.Request, res *kornet.Response) {

  c.kReq = req
  c.kRes = res

  if req.Body.Size() > 0 {

    c.body = req.Body.ReadAll()
    req.Body.Seek(0)
  }
}
