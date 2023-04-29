package app

import (
	"skfw/papaya"
	"skfw/papaya/pigeon/cors"
)

func ManageControlResourceShared(pn papaya.NetImpl) error {

	manageConsumers, err := cors.ManageConsumersNew()
	if err != nil {

		return err
	}

	// example using cors by consumer
	manageConsumers.Add("GET", "https://google.com")
	manageConsumers.GrantAll("https://google.com")

	pn.Use(cors.MakeMiddlewareForManageConsumers(manageConsumers))

	return nil
}
