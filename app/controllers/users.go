package controllers

import (
	"PapayaNet/papaya/bunny/swag"
	m "PapayaNet/papaya/koala/mapping"
	"net/http"
)

type Say struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func UserController(router swag.SwagRouterImpl) {

	router.Get("",
		&m.KMap{
			"permit":      true,
			"description": "Catch All Users",
			"params": &m.KMap{
				"q": "string",
			},
			"headers": &m.KMap{
				"auth": "string",
			},
			"request": &m.KMap{
				"application/json": &m.KMap{
					"description": "OK",
					"schema": &m.KMap{
						"name": "string",
						"info": &m.KMap{
							"r": []bool{},
							"metadata": []m.KMap{
								{
									"name": "string",
								},
							},
							"v":   "number",
							"say": &Say{},
						},
					},
				},
			},
			"responses": &m.KMap{
				"200": &m.KMap{
					"application/json": &m.KMap{
						"description": "OK",
						"schema": &m.KMap{
							"message": "string",
						},
					},
				},
			},
		},
		func(ctx *swag.SwagContext) error {

			ctx.Status(http.StatusOK)
			return ctx.JSON(m.KMap{
				"message": "Hello, World!",
			})
		})
}
