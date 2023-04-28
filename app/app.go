package app

import (
	"skfw/app/controllers"
	"skfw/papaya"
	"skfw/papaya/bunny/swag"
	"skfw/papaya/pigeon/templates/basicAuth/models"
	"skfw/papaya/pigeon/templates/basicAuth/repository"
	"time"
)

func App(pn papaya.NetImpl) error {

	var err error

	conn := pn.Connection()
	gorm := conn.GORM()

	swagger := pn.MakeSwagger(&swag.SwagInfo{
		Title:       "Example API",
		Version:     "1.0.0",
		Description: "Example API for documentation",
	})

	mainGroup := swagger.Group("/api/v1", "Schema")
	userGroup := mainGroup.Group("/users", "Authentication")

	userRouter := userGroup.Router()

	expired := time.Hour * 4
	activeDuration := time.Minute * 30 // interval
	maxSessions := 6

	basicAuth := repository.BasicAuthNew(conn, expired, activeDuration, maxSessions)
	basicAuth.Bind(swagger, userRouter)

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
