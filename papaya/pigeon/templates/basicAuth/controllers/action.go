package controllers

import (
  "bytes"
  "encoding/json"
  "errors"
  "github.com/golang-jwt/jwt/v5"
  "net/http"
  "skfw/papaya/bunny/swag"
  "skfw/papaya/koala/kornet"
  m "skfw/papaya/koala/mapping"
  "skfw/papaya/pigeon/drivers/common"
  "skfw/papaya/pigeon/drivers/postgresql"
  "skfw/papaya/pigeon/templates/basicAuth/models"
  "skfw/papaya/pigeon/templates/basicAuth/repository"
  "skfw/papaya/pigeon/templates/basicAuth/util"
  "time"

  "github.com/google/uuid"
  "gorm.io/gorm"
)

type BasicAuth struct {
  conn           common.DBConnectionImpl
  gorm           *gorm.DB
  userRepo       repository.UserRepositoryImpl
  sessionRepo    repository.SessionRepositoryImpl
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

  s.userRepo, err = repository.UserRepositoryNew(s.gorm)

  if err != nil {

    return err
  }

  s.sessionRepo, err = repository.SessionRepositoryNew(s.gorm)

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

    var err error
    var found bool

    var user *models.UserModel
    var session *models.SessionModel

    var claims jwt.MapClaims

    token := util.RequestAuth(ctx.Request())
    currentTime := time.Now().UTC()

    if token != "" {

      hashToken := util.HashSHA3(token)

      if session, found = s.sessionRepo.SearchFast(uuid.UUID{}, hashToken); found {

        //////// Recovery Token By Database ////////

        // max Session + 1
        // cause this method just validation token, not create another token

        if err = s.sessionRepo.RecoveryFast(util.Ids(session.UserID), hashToken, s.activeDuration, s.maxSessions+1); err != nil {

          switch err {
          case repository.TokenExpiredOrUserNoLongerActive, repository.SessionReachedLimit:

            ctx.Prevent()
            return ctx.Status(http.StatusUnauthorized).JSON(kornet.ResultNew(kornet.MessageNew(err.Error(), true), nil))
          }

          ctx.Prevent()
          return ctx.Status(http.StatusInternalServerError).JSON(kornet.ResultNew(kornet.MessageNew(err.Error(), true), nil))
        }

        //////// Recovery Token By Database ////////

        //////// JWT Checker ////////

        if claims, err = util.DecodeJWT(token, session.SecretKey); err != nil {

          ctx.Prevent()
          return ctx.Status(http.StatusUnauthorized).JSON(kornet.ResultNew(kornet.MessageNew("Invalid JWT", true), nil))
        }

        if uid, ok := claims["uid"]; ok {

          userId := m.KValueToString(uid)

          if !bytes.Equal([]byte(session.UserID), []byte(userId)) {

            ctx.Prevent()
            return ctx.Status(http.StatusUnauthorized).JSON(kornet.ResultNew(kornet.MessageNew("Invalid JWT", true), nil))
          }

        } else {

          ctx.Prevent()
          return ctx.Status(http.StatusUnauthorized).JSON(kornet.ResultNew(kornet.MessageNew("Invalid JWT", true), nil))
        }

        //////// JWT Checker ////////

        if util.DeviceRecognition(ctx, session) {

          if currentTime.Before(session.Expired) {

            if session.LastActivated.Before(currentTime.Add(s.activeDuration)) {

              if user, found = s.userRepo.SearchFastById(util.Ids(session.UserID)); found {

                if err = s.sessionRepo.CheckIn(session); err != nil {

                  ctx.Prevent()
                  return ctx.Status(http.StatusInternalServerError).JSON(kornet.ResultNew(kornet.MessageNew(err.Error(), true), nil))
                }

                ctx.Solve(user)
                return nil
              }

              ctx.Prevent()
              return ctx.Status(http.StatusInternalServerError).JSON(kornet.ResultNew(kornet.MessageNew("unable to get user information", true), nil))
            }

            ctx.Prevent()
            return ctx.Status(http.StatusUnauthorized).JSON(kornet.ResultNew(kornet.MessageNew("user is no longer active", true), nil))
          }

          ctx.Prevent()
          return ctx.Status(http.StatusUnauthorized).JSON(kornet.ResultNew(kornet.MessageNew("Token has expired", true), nil))
        }

        ctx.Prevent()
        return ctx.Status(http.StatusUnauthorized).JSON(kornet.ResultNew(kornet.MessageNew("Device not recognized", true), nil))
      }

      ctx.Prevent()
      return ctx.Status(http.StatusUnauthorized).JSON(kornet.ResultNew(kornet.MessageNew("Token not recognized", true), nil))
    }

    ctx.Prevent()
    return ctx.Status(http.StatusUnauthorized).JSON(kornet.ResultNew(kornet.MessageNew("Unauthorized", true), nil))
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
      "responses": swag.OkJSON(&kornet.Result{}),
    },
    func(ctx *swag.SwagContext) error {

      if ctx.Event() {

        req, _ := ctx.Kornet()

        sessionId, err := uuid.Parse(m.KValueToString(req.Query.Get("id")))

        if err != nil {

          return ctx.Status(http.StatusBadRequest).JSON(kornet.ResultNew(kornet.MessageNew("unable to parse uuid", true), nil))
        }

        if !util.EmptyIdx(sessionId) {

          if user, ok := ctx.Target().(*models.UserModel); ok {

            if session, ok := s.sessionRepo.SearchFastById(sessionId); ok {

              if session.UserID == user.ID {

                err := s.sessionRepo.Delete(session)
                if err != nil {

                  return ctx.Status(http.StatusInternalServerError).JSON(kornet.ResultNew(kornet.MessageNew("unable to remove session", true), nil))
                }

                return ctx.Status(http.StatusOK).JSON(kornet.ResultNew(kornet.MessageNew("remove session", false), nil))
              }

              return ctx.Status(http.StatusUnauthorized).JSON(kornet.ResultNew(kornet.MessageNew("Unauthorized", true), nil))
            }

            return ctx.Status(http.StatusInternalServerError).JSON(kornet.ResultNew(kornet.MessageNew("unable to get session information", true), nil))
          }

          return ctx.Status(http.StatusInternalServerError).JSON(kornet.ResultNew(kornet.MessageNew("unable to get user information", true), nil))
        }

        return ctx.Status(http.StatusBadRequest).JSON(kornet.ResultNew(kornet.MessageNew("uuid is empty", true), nil))
      }

      return ctx.Status(http.StatusNotAcceptable).JSON(kornet.ResultNew(kornet.MessageNew("not accepted by the event", true), nil))
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
      "responses": swag.OkJSON(&kornet.Result{}),
    },
    func(ctx *swag.SwagContext) error {

      if ctx.Event() {

        if user, ok := ctx.Target().(*models.UserModel); ok {

          if err := s.sessionRepo.DeleteFast(util.Ids(user.ID), "*"); err != nil {

            return ctx.Status(http.StatusInternalServerError).JSON(kornet.ResultNew(kornet.MessageNew("unable to remove all sessions", true), nil))
          }

          return ctx.Status(http.StatusOK).JSON(kornet.ResultNew(kornet.MessageNew("remove all sessions", false), nil))
        }
      }

      return ctx.Status(http.StatusNotAcceptable).JSON(kornet.ResultNew(kornet.MessageNew("not accepted by the event", true), nil))
    })

  router.Get("/sessions",
    &m.KMap{
      "AuthToken":   true,
      "description": "Catch All Sessions",
      "request": &m.KMap{
        "headers": &m.KMap{
          "Authorization": "string",
        },
      },
      "responses": swag.OkJSON(&kornet.Result{
        Data: []m.KMap{
          {
            "id":             "string",
            "used":           "boolean",
            "client_ip":      "string",
            "user_agent":     "string",
            "last_activated": "number",
            "expired":        "number",
          },
        },
      }),
    },
    func(ctx *swag.SwagContext) error {

      if ctx.Event() {

        if user, ok := ctx.Target().(*models.UserModel); ok {

          token := util.RequestAuth(ctx.Request())

          sessions, err := s.sessionRepo.GetAll(util.Ids(user.ID), s.maxSessions)
          if err != nil {

            return ctx.Status(http.StatusInternalServerError).JSON(kornet.ResultNew(kornet.MessageNew("unable to delete all sessions", true), nil))
          }

          var res []m.KMap
          res = make([]m.KMap, 0)

          // normalize
          for _, session := range sessions {

            res = append(res, m.KMap{

              "id":             session.ID,
              "used":           util.HashCompareSHA3(token, session.Token),
              "client_ip":      session.ClientIP,
              "user_agent":     session.UserAgent,
              "last_activated": session.LastActivated.UnixMilli(),
              "expired":        session.Expired.UnixMilli(),
            })
          }

          return ctx.Status(http.StatusOK).JSON(kornet.ResultNew(kornet.MessageNew("get all sessions", false), res))
        }
      }

      return ctx.Status(http.StatusNotAcceptable).JSON(kornet.ResultNew(kornet.MessageNew("not accepted by the event", true), nil))
    })
}

func (s *BasicAuth) MakeUserLoginEndpoint(router swag.SwagRouterImpl) {

  router.Post("/login",
    &m.KMap{
      "description": "Login",
      "request": &m.KMap{
        "headers": &m.KMap{
          "Authorization":    "string",
          "X-Forwarded-For?": "string",
          "X-Real-IP?":       "string",
        },
        "body": &m.KMap{
          "application/json": &m.KMap{
            "schema": &m.KMap{
              "username": "string",
              "email":    "string",
              "password": "string",
            },
          },
        },
      },
      "responses": swag.CreatedJSON(&kornet.Result{
        Data: &m.KMap{
          "token": "string",
        },
      }),
    },
    func(ctx *swag.SwagContext) error {

      ClientIP := util.GetClientIP(ctx)

      buff := ctx.Body()

      var err error

      var mm map[string]any

      mm = map[string]any{}

      err = json.Unmarshal(buff, &mm)
      if err != nil {

        return ctx.Status(http.StatusInternalServerError).JSON(kornet.ResultNew(kornet.MessageNew("unable to parse body into json", true), nil))
      }

      var ok, found bool
      var username, email, password string
      var user *models.UserModel
      var token string

      currentTime := time.Now().UTC()

      data := m.KMap(mm)

      username = m.KValueToString(data.Get("username"))
      email = m.KValueToString(data.Get("email"))
      password = m.KValueToString(data.Get("password"))

      if user, ok = s.userRepo.SearchFast(username, email); ok {

        if err = s.sessionRepo.RecoveryFast(util.Ids(user.ID), "", s.activeDuration, s.maxSessions); err != nil {

          if err != repository.TokenExpiredOrUserNoLongerActive {

            ctx.Prevent()
            return ctx.Status(http.StatusUnauthorized).JSON(kornet.ResultNew(kornet.MessageNew(err.Error(), true), nil))
          }
        }

        if user, found = s.userRepo.SearchFast(username, email); found {

          if util.CheckPasswordHash(password, user.Password) {

            clientIP := ClientIP
            userAgent := ctx.Get("User-Agent")

            token, err = s.sessionRepo.CreateFastAutoToken(user, clientIP, userAgent, currentTime.Add(s.expired), s.activeDuration, s.maxSessions)
            if err != nil {

              return ctx.Status(http.StatusInternalServerError).JSON(kornet.ResultNew(kornet.MessageNew(err.Error(), true), nil))
            }

            return ctx.Status(http.StatusCreated).JSON(kornet.ResultNew(kornet.MessageNew("login successful", false), &m.KMap{
              "token": token,
            }))
          }

          return ctx.Status(http.StatusUnauthorized).JSON(kornet.ResultNew(kornet.MessageNew("wrong password", true), nil))
        }
      }

      return ctx.Status(http.StatusUnauthorized).JSON(kornet.ResultNew(kornet.MessageNew("username, or email not found", true), nil))
    })
}

func (s *BasicAuth) MakeUserRegisterEndpoint(router swag.SwagRouterImpl) {

  router.Post("/signup",
    &m.KMap{
      "description": "Register New User",
      "request": &m.KMap{
        "body": swag.JSON(&m.KMap{
          "name":     "string", // real name | full name
          "username": "string",
          "email":    "string",
          "password": "string",
        }),
      },
      "responses": swag.CreatedJSON(&kornet.Result{
        Data: &m.KMap{
          "token": "string",
        },
      }),
    },
    func(ctx *swag.SwagContext) error {

      buff := ctx.Body()

      var err error

      var mm map[string]any

      mm = map[string]any{}

      err = json.Unmarshal(buff, &mm)
      if err != nil {

        return ctx.Status(http.StatusInternalServerError).JSON(kornet.ResultNew(kornet.MessageNew("unable to parse body into json", true), nil))
      }

      var user *models.UserModel
      var token string

      currentTime := time.Now().UTC()

      data := m.KMap(mm)

      var name, username string

      name = m.KValueToString(data.Get("name"))
      username = m.KValueToString(data.Get("username"))

      var valid bool
      var email util.EmailImpl
      var password util.PasswordImpl

      email, err = util.EmailNew(m.KValueToString(data.Get("email")))
      if err != nil {

        return ctx.Status(http.StatusBadRequest).JSON(kornet.ResultNew(kornet.MessageNew(err.Error(), true), nil))
      }

      if valid, err = email.Verify(); !valid {

        return ctx.Status(http.StatusBadRequest).JSON(kornet.ResultNew(kornet.MessageNew(err.Error(), true), nil))
      }

      password, err = util.PasswordNew(m.KValueToString(data.Get("password")))
      if err != nil {

        return ctx.Status(http.StatusBadRequest).JSON(kornet.ResultNew(kornet.MessageNew(err.Error(), true), nil))
      }

      // password contains
      minChars := 8
      specialChar := true
      numberChar := true
      upperChar := true
      lowerChar := true

      if valid, err = password.Verify(minChars, specialChar, numberChar, upperChar, lowerChar); !valid {

        return ctx.Status(http.StatusBadRequest).JSON(kornet.ResultNew(kornet.MessageNew(err.Error(), true), nil))
      }

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

      user, err = s.userRepo.CreateFast(name, username, email.Value(), password.Value()) // auto hashing password
      if err != nil {

        return ctx.Status(http.StatusInternalServerError).JSON(kornet.ResultNew(kornet.MessageNew(err.Error(), true), nil))
      }

      token, err = s.sessionRepo.CreateFastAutoToken(user, ctx.IP(), ctx.Get("User-Agent"), currentTime.Add(s.expired), s.activeDuration, s.maxSessions)
      if err != nil {

        return ctx.Status(http.StatusInternalServerError).JSON(kornet.ResultNew(kornet.MessageNew(err.Error(), true), nil))
      }

      return ctx.Status(http.StatusCreated).JSON(kornet.ResultNew(kornet.MessageNew("create new user successful", false), &m.KMap{
        "token": token,
      }))
    })
}
