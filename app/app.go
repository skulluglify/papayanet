package app

import (
	"skfw/app/controllers"
	"skfw/papaya"
	"skfw/papaya/bunny/swag"
	"skfw/papaya/pigeon/templates/basic/repository"
)

func App(pn papaya.NetImpl) error {

	swagger := pn.MakeSwagger(&swag.SwagInfo{
		Title:       "Example API",
		Version:     "1.0.0",
		Description: "Example API for documentation",
	})

	simpleAction := repository.SimpleActionNew(pn.Connection())
	swagger.AddTask(simpleAction.MakeAuthTokenTask())

	mainGroup := swagger.Group("/api/v1", "Schema")
	userGroup := mainGroup.Group("/users", "Authentication")

	controllers.UserController(pn, userGroup.Router())

	if err := swagger.Start(); err != nil {
		return err
	}

	return pn.Serve("127.0.0.1", 8000)
}
