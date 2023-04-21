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

	router.Get("/:id",
		&m.KMap{
			"permit":      true,
			"description": "Catch All Users",
			"request": &m.KMap{
				"task": true,
				"params": &m.KMap{
					"#id": "number",
					"q":   "number",
				},
				"headers": &m.KMap{
					"auth": "string",
				},
			},
			"responses": swag.OkJSON(Say{}),
		},
		func(ctx *swag.SwagContext) error {

			ctx.Status(http.StatusOK)
			return ctx.JSON(Say{200, "OK"})
		})

	router.Post("/:id",
		&m.KMap{
			"permit":      true,
			"description": "Catch All Users",
			"request": &m.KMap{
				"params": &m.KMap{
					"#id": "number",
					"q":   "string",
				},
				"headers": &m.KMap{
					"auth": "string",
				},
				"body": swag.JSON(&m.KMap{
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
				}),
			},
			"responses": swag.OkJSON(Say{}),
		},
		func(ctx *swag.SwagContext) error {

			ctx.Status(http.StatusOK)
			return ctx.JSON(m.KMap{
				"message": "Hello, World!",
			})
		})
}
