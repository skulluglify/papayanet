package papaya

import (
  "os"
  "skfw/papaya/bunny/swag"
  "skfw/papaya/koala"
  m "skfw/papaya/koala/mapping"
  "skfw/papaya/pigeon"
  "skfw/papaya/pigeon/drivers/common"
  "skfw/papaya/pigeon/drivers/postgresql"
  "skfw/papaya/util"
  "strconv"

  "github.com/gofiber/fiber/v2"
  "github.com/joho/godotenv"
)

type Net struct {
  Console koala.KConsoleImpl
  *fiber.Config
  *fiber.App
  *postgresql.DBConfig
  DBConnection common.DBConnectionImpl
  version      koala.KVersionImpl
}

type NetImpl interface {
  Init()
  Serve(host string, port int) error
  MakeSwagger(info *swag.SwagInfo) swag.SwagImpl
  Logger() koala.KConsoleImpl
  Connection() common.DBConnectionImpl
  Version() koala.KVersionImpl
  Use(args ...any)
  Close() error
}

func NetNew() NetImpl {

  net := &Net{
    Config: &fiber.Config{
      DisableStartupMessage: true,
    },
    DBConfig: &postgresql.DBConfig{},
  }

  net.Init()

  return net
}

func (n *Net) Init() {

  if n.Console == nil {

    n.Console = koala.KConsoleNew()
  }

  // Load `.env`
  if err := godotenv.Load(); err != nil {

    n.Console.Error(err)
  }

  if n.App == nil {

    n.Config.ErrorHandler = func(ctx *fiber.Ctx, err error) error {

      return ctx.Status(fiber.StatusNotFound).JSON(&m.KMap{
        "message": "site not found",
        "error":   true,
      })
    }

    n.App = fiber.New(*n.Config)
  }

  if n.DBConnection == nil {

    var err error
    n.DBConnection, err = postgresql.DBConnectionNew(pigeon.InitLoadEnviron)

    if err != nil {

      n.Console.Error(err)
      os.Exit(1)
    }
  }

  n.version = koala.KVersionNew(
    util.VersionMajor,
    util.VersionMinor,
    util.VersionPatch,
  )

  n.Console.Log(n.Console.Text(util.Banner(n.version), koala.ColorGreen, koala.ColorBlack, koala.StyleBold))
  n.Console.Log("Server has started ...")
}

func (n *Net) Serve(host string, port int) error {

  return n.App.Listen(host + ":" + strconv.Itoa(port))
}

func (n *Net) MakeSwagger(info *swag.SwagInfo) swag.SwagImpl {

  return swag.MakeSwag(n.App, info)
}

func (n *Net) Connection() common.DBConnectionImpl {

  return n.DBConnection
}

func (n *Net) Version() koala.KVersionImpl {

  return n.version
}

func (n *Net) Logger() koala.KConsoleImpl {

  return n.Console
}

func (n *Net) Use(args ...any) {

  if n.App != nil {

    // registry new middleware
    n.App.Use(args...)
  }
}

func (n *Net) Close() error {

  if err := n.App.Shutdown(); err != nil {

    return err
  }

  if err := n.DBConnection.Close(); err != nil {

    return err
  }

  n.Console.Log("Server has shutdown ...")

  return nil
}
