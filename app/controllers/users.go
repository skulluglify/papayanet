package controllers

import (
	"encoding/json"
	"net/http"
	"skfw/papaya"
	"skfw/papaya/bunny/swag"
	"skfw/papaya/koala/kornet"
	m "skfw/papaya/koala/mapping"
	"skfw/papaya/pigeon/templates/basic/models"
	"skfw/papaya/pigeon/templates/basic/repository"
	"time"
)

func UserController(pn papaya.NetImpl, router swag.SwagRouterImpl) {

	conn := pn.Connection()
	gorm := conn.GORM()

	userRepo := repository.UserRepositoryNew(gorm)
	sessionRepo := repository.SessionRepositoryNew(gorm)

	router.Get("/test",
		&m.KMap{
			"AuthToken":   true,
			"description": "Test Authentication",
			"request": &m.KMap{
				"headers": &m.KMap{
					"Authorization": "string",
				},
			},
			"responses": swag.OkJSON(&models.UserModel{}),
		},
		func(ctx *swag.SwagContext) error {

			if user, ok := ctx.Event().(models.UserModel); ok {

				ctx.Status(http.StatusOK)
				return ctx.JSON(user)
			}

			return ctx.JSON(kornet.MessageNew("can't catch user data", true))
		})

	router.Post("/login",
		&m.KMap{
			"description": "Login",
			"request": &m.KMap{
				"headers": &m.KMap{
					"Authorization": "string",
				},
				"body": swag.JSON(m.KMap{
					"username": "string",
					"email":    "string",
					"password": "string",
				}),
			},
			"responses": swag.OkJSON(models.SessionModel{}),
		},
		func(ctx *swag.SwagContext) error {

			return nil
		})

	router.Post("/signup",
		&m.KMap{
			"description": "Register New User",
			"request": &m.KMap{
				"body": swag.JSON(m.KMap{
					"username": "string",
					"email":    "string",
					"password": "string",
				}),
			},
			"responses": swag.CreatedJSON(models.SessionModel{}),
		},
		func(ctx *swag.SwagContext) error {

			buff := ctx.Body()

			var err error

			var data map[string]any

			data = map[string]any{}

			err = json.Unmarshal(buff, &data)
			if err != nil {

				return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to parse body into json", true))
			}

			var ok bool
			var username, email, password string
			var user *models.UserModel

			var session *models.SessionModel

			if username, ok = data["username"].(string); ok {
				if email, ok = data["email"].(string); ok {
					if password, ok = data["password"].(string); ok {

						user, err = userRepo.CreateFast(username, email, password)
						if err != nil {

							return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
						}

						session, err = sessionRepo.CreateFastAutoToken(user, time.Now().UTC().Add(time.Hour*4))
						if err != nil {

							return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
						}

						return ctx.Status(http.StatusCreated).JSON(session)
					}
				}
			}

			return nil
		})
}
