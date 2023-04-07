package app

import (
	"PapayaNet/app/controllers"
	"PapayaNet/papaya"
	"PapayaNet/papaya/swag"
)

func App(pn papaya.NetImpl) error {

	swagInfo := swag.MakeSwagInfo("Example API", "Example API for documentation")
	swagger := pn.MakeSwagger(swagInfo)

	router := swagger.Router()

	controllers.User(router)

	return pn.Serve("127.0.0.1", 8000)
}
