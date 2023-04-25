package swag

import "github.com/gofiber/fiber/v2"

type SwagContext struct {
	*fiber.Ctx
	event    bool
	data     any
	finished bool
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

func (c *SwagContext) Dispatch(data any) {

	c.data = data
}

func (c *SwagContext) Prevent() {

	c.finished = true
}

func (c *SwagContext) Override() bool {

	return c.finished
}
