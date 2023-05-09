package app

import (
	"skfw/app/controllers"
	"skfw/papaya"
	"skfw/papaya/bunny/swag"
	bac "skfw/papaya/pigeon/templates/basicAuth/controllers" // basic auth action
	"skfw/papaya/pigeon/templates/basicAuth/models"
	"time"
)

func App(pn papaya.NetImpl) error {

	var err error

	conn := pn.Connection()
	gorm := conn.GORM()

	ManageControlResourceShared(pn)

	swagger := pn.MakeSwagger(&swag.SwagInfo{
		Title:       "Example API",
		Version:     "1.0.0",
		Description: "Example API for Documentation",
	})

	mainGroup := swagger.Group("/api/v1", "Schema")
	userGroup := mainGroup.Group("/users", "Authentication")

	userRouter := userGroup.Router()

	expired := time.Hour * 4
	activeDuration := time.Minute * 30 // interval
	maxSessions := 6

	basicAuth := bac.BasicAuthNew(conn, expired, activeDuration, maxSessions)
	basicAuth.Bind(swagger, userRouter)

	// ref: https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
	//gorm.Exec("ALTER TABLE carts ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE")

	if err = gorm.AutoMigrate(&models.UserModel{}, &models.SessionModel{}); err != nil {

		return err
	}

	if err = controllers.UserController(pn, userRouter); err != nil {
		return err
	}

	if err = swagger.Start(); err != nil {
		return err
	}

	return pn.Serve("127.0.0.1", 8000)
}
