package app

import (
	"PapayaNet/app/controllers"
	"PapayaNet/papaya"
	"PapayaNet/papaya/bunny/swag"
	m "PapayaNet/papaya/koala/mapping"
)

func App(pn papaya.NetImpl) error {

	swagger := pn.MakeSwagger(&swag.SwagInfo{
		Title:       "Example API",
		Version:     "1.0.0",
		Description: "Example API for documentation",
	})

	swagger.AddTask("request.task", func(ctx *swag.SwagContext) error {

		ctx.Closed = true

		pn.Logger().Log("need request.task", ctx.Event())

		return ctx.JSON(m.KMap{
			"message": "task running",
		})
	})

	mainGroup := swagger.Group("/api/v1", "Schema")
	userGroup := mainGroup.Group("/users", "Authentication")

	controllers.UserController(userGroup.Router())

	if err := swagger.Start(); err != nil {
		return err
	}

	return pn.Serve("127.0.0.1", 8000)
}
