package app

import (
	"PapayaNet/app/controllers"
	"PapayaNet/papaya"
	"PapayaNet/papaya/bunny/swag"
	"PapayaNet/papaya/pigeon/templates/basic"
)

func App(pn papaya.NetImpl) error {

	swagger := pn.MakeSwagger(&swag.SwagInfo{
		Title:       "Example API",
		Version:     "1.0.0",
		Description: "Example API for documentation",
	})

	//conn := pn.Connection()
	//db, _ := conn.Database()

	swagger.AddTask(swag.MakeSwagTask("request.task", func(ctx *swag.SwagContext) error {

		//ctx.Prevent()

		pn.Logger().Log("task running 1 ...", ctx.Event())

		return nil
	}))

	swagger.AddTask(swag.MakeSwagTask("request.permit", func(ctx *swag.SwagContext) error {

		pn.Logger().Log("task permit granted ...", ctx.Event())

		return nil
	}))

	swagger.AddTask(swag.MakeSwagTask("request.task", func(ctx *swag.SwagContext) error {

		pn.Logger().Log("task running 2 ...", ctx.Event())

		return nil
	}))

	swagger.AddTask(basic.MakeAuthTokenTask(pn))

	mainGroup := swagger.Group("/api/v1", "Schema")
	userGroup := mainGroup.Group("/users", "Authentication")

	controllers.UserController(userGroup.Router())

	if err := swagger.Start(); err != nil {
		return err
	}

	return pn.Serve("127.0.0.1", 8000)
}
