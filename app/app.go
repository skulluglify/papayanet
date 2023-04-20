package app

import (
	"PapayaNet/app/controllers"
	"PapayaNet/papaya"
	"PapayaNet/papaya/bunny/swag"
)

func App(pn papaya.NetImpl) error {

	swagger := pn.MakeSwagger(&swag.SwagInfo{
		Title:       "Example API",
		Version:     "1.0.0",
		Description: "Example API for documentation",
	})

	mainGroup := swagger.Group("/api/v1", "Schema")
	userGroup := mainGroup.Group("/users", "Authentication")

	controllers.UserController(userGroup.Router())

	if err := swagger.Start(); err != nil {
		return err
	}

	swagger.Swagger()

	return pn.Serve("127.0.0.1", 8000)
	// return nil
}
