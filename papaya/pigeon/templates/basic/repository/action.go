package repository

import (
  "errors"
  "github.com/golang-jwt/jwt/v5"
  "gorm.io/gorm"
  "net/http"
  "skfw/papaya/bunny/swag"
  "skfw/papaya/koala/kornet"
  m "skfw/papaya/koala/mapping"
  "skfw/papaya/pigeon/drivers/postgresql"
  "skfw/papaya/pigeon/templates/basic/models"
  "time"
)

type SimpleAction struct {
  conn        postgresql.DBConnectionImpl
  gorm        *gorm.DB
  userRepo    UserRepositoryImpl
  sessionRepo SessionRepositoryImpl
}

type SimpleActionImpl interface {
  Init(conn postgresql.DBConnectionImpl) error
  MakeAuthTokenTask() *swag.SwagTask
}

func SimpleActionNew(conn postgresql.DBConnectionImpl) SimpleActionImpl {

  simpleAction := &SimpleAction{}
  err := simpleAction.Init(conn)
  if err != nil {
    return nil
  }
  return simpleAction
}

func (s *SimpleAction) Init(conn postgresql.DBConnectionImpl) error {

  var err error

  if conn == nil {

    return errors.New("conn is NULL")
  }

  s.conn = conn
  s.gorm = conn.GORM()

  user := &models.UserModel{}
  session := &models.SessionModel{}

  s.userRepo = UserRepositoryNew(s.gorm)
  s.sessionRepo = SessionRepositoryNew(s.gorm)

  // check null
  if s.userRepo == nil || s.sessionRepo == nil {

    return nil
  }

  // enable extension UUID
  err = postgresql.PgEnableExtensionUUID(conn)
  if err != nil {
    return errors.New("failed to enable extension UUID")
  }

  // set time zone as UTC
  err = postgresql.PgSetTimeZoneUTC(conn)
  if err != nil {
    return errors.New("failed to set TimeZone as UTC")
  }

  // auto migration
  err = s.gorm.AutoMigrate(user, session)
  if err != nil {
    return errors.New("failed to migrate database")
  }

  return nil
}

func (s *SimpleAction) MakeAuthTokenTask() *swag.SwagTask {

  // make auth task
  return swag.MakeSwagTask("AuthToken", func(ctx *swag.SwagContext) error {

    if m.KValueToBool(ctx.Event()) {

      var ok, found bool
      var username, email string
      var claims jwt.MapClaims
      var data map[string]any

      var user *models.UserModel
      var session *models.SessionModel

      req := ctx.Request()

      auth := RequestAuth(req)

      if auth != "" {

        claims, _ = DecodeJWT(auth, "", time.Time{})
        data = claims

        if username, ok = data["username"].(string); ok {

          if email, ok = data["email"].(string); ok {

            if user, found = s.userRepo.SearchFast(username, email); found {

              if session, found = s.sessionRepo.SearchFast(user.ID, auth); found {

                if session.Expired.Before(time.Now().UTC()) {

                  ctx.Dispatch(user)
                  return nil
                }

                ctx.Prevent()
                return ctx.Status(http.StatusUnauthorized).JSON(kornet.MessageNew("Token has expired", true))
              }
            }
          }
        }
      }
    }

    ctx.Prevent()
    return ctx.Status(http.StatusUnauthorized).JSON(kornet.MessageNew("Unauthorized", true))
  })
}
