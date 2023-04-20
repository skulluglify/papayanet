package papaya

import (
  swag2 "PapayaNet/papaya/bunny/swag"
  "PapayaNet/papaya/koala"
  m "PapayaNet/papaya/koala/mapping"
  "PapayaNet/papaya/pigeon"
  "PapayaNet/papaya/pigeon/drivers/postgresql"
  "PapayaNet/papaya/util"
  "os"
  "strconv"

  "github.com/gofiber/fiber/v2"
  "github.com/joho/godotenv"
)

type Net struct {
  Console koala.KConsoleImpl
  *fiber.Config
  *fiber.App
  *postgresql.DBConfig
  *postgresql.DBConnection
  Version koala.KVersionImpl
}

type NetImpl interface {
  Init()
  Serve(host string, port int) error
  MakeSwagger(info *swag2.SwagInfo) swag2.SwagImpl
  Logger() koala.KConsoleImpl
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

  n.Version = koala.KVersionNew(
    util.VersionMajor,
    util.VersionMinor,
    util.VersionPatch,
  )

  n.Console.Log(n.Console.Text(util.Banner(n.Version), koala.ColorGreen, koala.ColorBlack, koala.StyleBold))
  n.Console.Log("Server has started ...")
}

func (n *Net) Serve(host string, port int) error {

  return n.App.Listen(host + ":" + strconv.Itoa(port))
}

func (n *Net) MakeSwagger(info *swag2.SwagInfo) swag2.SwagImpl {

  return swag2.MakeSwag(n.App, info)
}

func (n *Net) Logger() koala.KConsoleImpl {

  return n.Console
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
