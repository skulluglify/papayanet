package repository

import (
  "encoding/json"
  "errors"
  "net/http"
  "skfw/papaya/bunny/swag"
  "skfw/papaya/koala/kornet"
  m "skfw/papaya/koala/mapping"
  "skfw/papaya/pigeon/drivers/common"
  "skfw/papaya/pigeon/drivers/postgresql"
  "skfw/papaya/pigeon/templates/basicAuth/models"
  "time"

  "github.com/google/uuid"
  "gorm.io/gorm"
)

type BasicAuth struct {
  conn           common.DBConnectionImpl
  gorm           *gorm.DB
  userRepo       UserRepositoryImpl
  sessionRepo    SessionRepositoryImpl
  expired        time.Duration
  activeDuration time.Duration
  maxSessions    int
}

type BasicAuthImpl interface {
  Init(conn common.DBConnectionImpl, expired time.Duration, activeDuration time.Duration, maxSessions int) error
  Bind(swag swag.SwagImpl, router swag.SwagRouterImpl)
  Migration() error

  MakeAuthTokenTask() *swag.SwagTask
  MakeSessionEndpoint(router swag.SwagRouterImpl)
  MakeUserLoginEndpoint(router swag.SwagRouterImpl)
  MakeUserRegisterEndpoint(router swag.SwagRouterImpl)
}

func BasicAuthNew(conn common.DBConnectionImpl, expired time.Duration, activeDuration time.Duration, maxSessions int) BasicAuthImpl {

  simpleAction := &BasicAuth{}
  err := simpleAction.Init(conn, expired, activeDuration, maxSessions)
  if err != nil {
    return nil
  }
  return simpleAction
}

func (s *BasicAuth) Init(conn common.DBConnectionImpl, expired time.Duration, activeDuration time.Duration, maxSessions int) error {

  var err error

  if conn == nil {

    return errors.New("conn is NULL")
  }

  s.conn = conn
  s.gorm = conn.GORM()

  s.userRepo, err = UserRepositoryNew(s.gorm)

  if err != nil {

    return err
  }

  s.sessionRepo, err = SessionRepositoryNew(s.gorm)

  if err != nil {

    return err
  }

  s.expired = expired
  s.activeDuration = activeDuration
  s.maxSessions = maxSessions

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

  if err != nil {
    return errors.New("failed to migrate database")
  }

  return nil
}

func (s *BasicAuth) Bind(swag swag.SwagImpl, router swag.SwagRouterImpl) {

  // register task
  swag.AddTask(s.MakeAuthTokenTask())

  // register endpoints
  s.MakeSessionEndpoint(router)
  s.MakeUserLoginEndpoint(router)
  s.MakeUserRegisterEndpoint(router)
}

func (s *BasicAuth) Migration() error {

  // auto migration
  if err := s.gorm.AutoMigrate(&models.UserModel{}, &models.SessionModel{}); err != nil {

    return err
  }

  return nil
}

func (s *BasicAuth) MakeAuthTokenTask() *swag.SwagTask {

  // make auth task
  return swag.MakeSwagTask("AuthToken", func(ctx *swag.SwagContext) error {

    var found bool

    var user *models.UserModel
    var session *models.SessionModel

    auth := RequestAuth(ctx.Request())
    currentTime := time.Now().UTC()

    if auth != "" {

      if session, found = s.sessionRepo.SearchFast(uuid.UUID{}, auth); found {

        if err := s.sessionRepo.RecoveryFast(Ids(session.UserID), auth, s.activeDuration, s.maxSessions); err != nil {

          switch err {
          case TokenExpiredOrUserNoLongerActive, SessionReachedLimit:

            ctx.Prevent()
            return ctx.Status(http.StatusUnauthorized).JSON(kornet.MessageNew(err.Error(), true))
          }

          ctx.Prevent()
          return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
        }

        if DeviceRecognition(session, ctx) {

          if currentTime.Before(session.Expired) {

            if session.LastActivated.Before(currentTime.Add(s.activeDuration)) {

              if user, found = s.userRepo.SearchFastById(Ids(session.UserID)); found {

                if err := s.sessionRepo.CheckIn(session); err != nil {

                  ctx.Prevent()
                  return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
                }

                ctx.Solve(user)
                return nil
              }

              ctx.Prevent()
              return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to get user information", true))
            }

            ctx.Prevent()
            return ctx.Status(http.StatusUnauthorized).JSON(kornet.MessageNew("user is no longer active", true))
          }

          ctx.Prevent()
          return ctx.Status(http.StatusUnauthorized).JSON(kornet.MessageNew("Token has expired", true))
        }

        ctx.Prevent()
        return ctx.Status(http.StatusUnauthorized).JSON(kornet.MessageNew("Device not recognized", true))
      }

      ctx.Prevent()
      return ctx.Status(http.StatusUnauthorized).JSON(kornet.MessageNew("Token not recognized", true))
    }

    ctx.Prevent()
    return ctx.Status(http.StatusUnauthorized).JSON(kornet.MessageNew("Unauthorized", true))
  })
}

func (s *BasicAuth) MakeSessionEndpoint(router swag.SwagRouterImpl) {

  router.Delete("/session",
    &m.KMap{
      "AuthToken":   true,
      "description": "Delete Session",
      "request": &m.KMap{
        "params": &m.KMap{
          "id": "string",
        },
        "headers": &m.KMap{
          "Authorization": "string",
        },
      },
      "responses": swag.OkJSON(kornet.Message{}),
    },
    func(ctx *swag.SwagContext) error {

      if ctx.Event() {

        req, _ := ctx.Kornet()

        sessionId, err := uuid.Parse(m.KValueToString(req.Query.Get("id")))

        if err != nil {

          return ctx.Status(http.StatusBadRequest).JSON(kornet.MessageNew("unable to parse uuid", true))
        }

        if !EmptyIdx(sessionId) {

          if user, ok := ctx.Target().(*models.UserModel); ok {

            if session, ok := s.sessionRepo.SearchFastById(sessionId); ok {

              if session.UserID == user.ID {

                err := s.sessionRepo.Delete(session)
                if err != nil {

                  return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to remove session", true))
                }

                return ctx.Status(http.StatusOK).JSON(kornet.MessageNew("remove session", false))
              }

              return ctx.Status(http.StatusUnauthorized).JSON(kornet.MessageNew("Unauthorized", true))
            }

            return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to get session information", true))
          }

          return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to get user information", true))
        }

        return ctx.Status(http.StatusBadRequest).JSON(kornet.MessageNew("uuid is empty", true))
      }

      return ctx.Status(http.StatusNotAcceptable).JSON(kornet.MessageNew("not accepted by the event", true))
    })

  router.Delete("/sessions",
    &m.KMap{
      "AuthToken":   true,
      "description": "Delete All Sessions",
      "request": &m.KMap{
        "headers": &m.KMap{
          "Authorization": "string",
        },
      },
      "responses": swag.OkJSON(kornet.Message{}),
    },
    func(ctx *swag.SwagContext) error {

      if ctx.Event() {

        if user, ok := ctx.Target().(*models.UserModel); ok {

          if err := s.sessionRepo.DeleteFast(Ids(user.ID), "*"); err != nil {

            return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to remove all sessions", true))
          }

          return ctx.Status(http.StatusOK).JSON(kornet.MessageNew("remove all sessions", false))
        }
      }

      return ctx.Status(http.StatusNotAcceptable).JSON(kornet.MessageNew("not accepted by the event", true))
    })

  router.Get("/sessions",
    &m.KMap{
      "AuthToken":   true,
      "description": "Delete All Sessions",
      "request": &m.KMap{
        "headers": &m.KMap{
          "Authorization": "string",
        },
      },
      "responses": swag.OkJSON([]m.KMap{
        {
          "id":             "string",
          "used":           "boolean",
          "client_ip":      "string",
          "user_agent":     "string",
          "last_activated": "number",
          "expired":        "number",
        },
      }),
    },
    func(ctx *swag.SwagContext) error {

      if ctx.Event() {

        if user, ok := ctx.Target().(*models.UserModel); ok {

          auth := RequestAuth(ctx.Request())

          sessions, err := s.sessionRepo.GetAll(Ids(user.ID), s.maxSessions)
          if err != nil {

            return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to delete all sessions", true))
          }

          var res []m.KMap
          res = make([]m.KMap, 0)

          // normalize
          for _, session := range sessions {

            res = append(res, m.KMap{

              "id":             session.ID,
              "used":           session.Token == auth,
              "client_ip":      session.ClientIP,
              "user_agent":     session.UserAgent,
              "last_activated": session.LastActivated.UnixMilli(),
              "expired":        session.Expired.UnixMilli(),
            })
          }

          return ctx.Status(http.StatusOK).JSON(res)
        }
      }

      return ctx.Status(http.StatusNotAcceptable).JSON(kornet.MessageNew("not accepted by the event", true))
    })
}

func (s *BasicAuth) MakeUserLoginEndpoint(router swag.SwagRouterImpl) {

  router.Post("/login",
    &m.KMap{
      "description": "Login",
      "request": &m.KMap{
        "headers": &m.KMap{
          "Authorization": "string",
        },
        "body": &m.KMap{
          "application/json": &m.KMap{
            "schema": &m.KMap{
              "username": "string",
              "email":    "string",
              "password": "string",
            },
          },
          "application/xml": &m.KMap{
            "schema": &m.KMap{
              "username": "string",
              "email":    "string",
              "password": "string",
            },
          },
        },
      },
      "responses": swag.CreatedJSON(&m.KMap{
        "token":   "string",
        "message": "string",
        "error":   "boolean",
      }),
    },
    func(ctx *swag.SwagContext) error {

      buff := ctx.Body()

      var err error

      var mm map[string]any

      mm = map[string]any{}

      err = json.Unmarshal(buff, &mm)
      if err != nil {

        return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to parse body into json", true))
      }

      var ok, found bool
      var username, email, password string
      var user *models.UserModel
      var session *models.SessionModel

      currentTime := time.Now().UTC()

      data := m.KMap(mm)

      username = m.KValueToString(data.Get("username"))
      email = m.KValueToString(data.Get("email"))
      password = m.KValueToString(data.Get("password"))

      if user, ok = s.userRepo.SearchFast(username, email); ok {

        if err = s.sessionRepo.RecoveryFast(Ids(user.ID), "", s.activeDuration, s.maxSessions); err != nil {

          if err != TokenExpiredOrUserNoLongerActive {

            ctx.Prevent()
            return ctx.Status(http.StatusUnauthorized).JSON(kornet.MessageNew(err.Error(), true))
          }
        }

        if user, found = s.userRepo.SearchFast(username, email); found {

          if CheckPasswordHash(password, user.Password) {

            clientIP := ctx.IP()
            userAgent := ctx.Get("User-Agent")

            session, err = s.sessionRepo.CreateFastAutoToken(user, clientIP, userAgent, currentTime.Add(s.expired), s.activeDuration, s.maxSessions)
            if err != nil {

              return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
            }

            return ctx.Status(http.StatusCreated).JSON(&m.KMap{
              "token":   session.Token,
              "message": "login successful",
              "error":   false,
            })
          }

          return ctx.Status(http.StatusUnauthorized).JSON(kornet.MessageNew("wrong password", true))
        }
      }

      return ctx.Status(http.StatusUnauthorized).JSON(kornet.MessageNew("username, or email not found", true))
    })
}

func (s *BasicAuth) MakeUserRegisterEndpoint(router swag.SwagRouterImpl) {

  router.Post("/signup",
    &m.KMap{
      "description": "Register New User",
      "request": &m.KMap{
        "body": swag.JSON(&m.KMap{
          "username": "string",
          "email":    "string",
          "password": "string",
        }),
      },
      "responses": swag.CreatedJSON(&m.KMap{
        "token":   "string",
        "message": "string",
        "error":   "boolean",
      }),
    },
    func(ctx *swag.SwagContext) error {

      buff := ctx.Body()

      var err error

      var mm map[string]any

      mm = map[string]any{}

      err = json.Unmarshal(buff, &mm)
      if err != nil {

        return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to parse body into json", true))
      }

      var user *models.UserModel
      var session *models.SessionModel

      currentTime := time.Now().UTC()

      data := m.KMap(mm)

      var username, email, password string

      username = m.KValueToString(data.Get("username"))
      email = m.KValueToString(data.Get("email"))
      password = m.KValueToString(data.Get("password"))

      //user = &models.UserModel{
      //  //Name:        m.KValueToString(data.Get("name")),
      //  Username: m.KValueToString(data.Get("username")),
      //  Email:    m.KValueToString(data.Get("email")),
      //  Password: m.KValueToString(data.Get("password")),
      //  //Gender:      m.KValueToString(data.Get("gender")),
      //  //Phone:       m.KValueToString(data.Get("phone")),
      //  //DOB:         time.UnixMilli(m.KValueToInt(data.Get("dob"))).UTC(), // make it relative to use in everywhere
      //  //Address:     m.KValueToString(data.Get("address")),
      //  //CountryCode: m.KValueToString(data.Get("country_code")),
      //  //City:        m.KValueToString(data.Get("city")),
      //  //PostalCode:  m.KValueToString(data.Get("postal_code")),
      //  Admin: false,
      //}

      //user.Password, err = HashPassword(user.Password)
      //if err != nil {
      //  return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to hash password", true))
      //}

      //err = s.userRepo.Create(user)
      user, err = s.userRepo.CreateFast(username, email, password) // auto hashing password
      if err != nil {

        return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
      }

      session, err = s.sessionRepo.CreateFastAutoToken(user, ctx.IP(), ctx.Get("User-Agent"), currentTime.Add(s.expired), s.activeDuration, s.maxSessions)
      if err != nil {

        return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
      }

      return ctx.Status(http.StatusCreated).JSON(m.KMap{
        "token":   session.Token,
        "message": "create new user successful",
        "error":   false,
      })
    })
}
