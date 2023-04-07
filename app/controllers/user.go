package controllers

import (
	"PapayaNet/papaya/koala/mapping"
	"PapayaNet/papaya/swag"
	"net/http"
	"reflect"
)

func User(router swag.SwagRouterImpl) {

	group := router.Group("user")
	router = group.Router()

	router.Get("/",
		&mapping.KMap{
			"params": &mapping.KMap{
				"q": reflect.String,
			},
			"headers": &mapping.KMap{
				"auth": reflect.String,
			},
			"request": &mapping.KMap{
				"application/json": &mapping.KMap{
					"name": reflect.String,
				},
			},
			"responses": &mapping.KMap{
				"200": &mapping.KMap{
					"application/json": &mapping.KMap{
						"message": reflect.String,
					},
				},
			},
		},
		func(ctx *swag.SwagContext) error {

			ctx.Status(http.StatusOK)
			return ctx.JSON(mapping.KMap{
				"message": "Hello, World!",
			})
		})
}
