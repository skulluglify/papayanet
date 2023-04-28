package app

import (
	"skfw/app/controllers"
	"skfw/papaya"
	"skfw/papaya/bunny/swag"
	"skfw/papaya/pigeon/templates/basicAuth/repository"
	"time"
)

func App(pn papaya.NetImpl) error {

	var err error

	swagger := pn.MakeSwagger(&swag.SwagInfo{
		Title:       "Example API",
		Version:     "1.0.0",
		Description: "Example API for documentation",
	})

	mainGroup := swagger.Group("/api/v1", "Schema")
	userGroup := mainGroup.Group("/users", "Authentication")

	userRouter := userGroup.Router()

	expired := time.Hour * 24
	activeDuration := time.Minute * 30
	maxSessions := 6

	basicAuth := repository.BasicAuthNew(pn.Connection(), expired, activeDuration, maxSessions)

	swagger.AddTask(basicAuth.MakeAuthTokenTask())

	basicAuth.MakeSessionEndpoint(userRouter)
	basicAuth.MakeUserLoginEndpoint(userRouter)
	basicAuth.MakeUserRegisterEndpoint(userRouter)

	err = controllers.UserController(pn, userRouter)
	if err != nil {
		return err
	}

	if err = swagger.Start(); err != nil {
		return err
	}

	return pn.Serve("127.0.0.1", 8000)
}
