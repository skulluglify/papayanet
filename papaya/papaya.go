package papaya

import (
	"PapayaNet/papaya/db"
	"PapayaNet/papaya/db/drivers/mysql"
	net2 "PapayaNet/papaya/peanut"
	"PapayaNet/papaya/utils"
	"crypto/tls"
	"errors"
	"github.com/labstack/echo/v4"
	gLog "github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strconv"
)

type PapayaNet struct {
	DB      *gorm.DB
	Echo    *echo.Echo
	Version *utils.PnVersion
	Console *utils.PnConsole
}

type PapayaNetImpl interface {
	Init() error
	EnvLoader(environ PnEnvImpl) error
	Serve(addr string, port uint16)
}

func (pn *PapayaNet) Init() error {

	if pn.Version == nil {

		pn.Version = utils.PnMakeVersion(
			PN_MAJOR_VERSION,
			PN_MINOR_VERSION,
			PN_PATCH_VERSION)
	}

	banner := Banner(pn.Version)

	if pn.Console == nil {

		console := utils.PnMakeConsole()
		//console.Colorful = true // force use colored
		//console.Silent = true // force silent mode
		console.Listen(func(info int) error {

			switch info {
			case utils.PN_CONSOLE_TYPE_ERROR, utils.PN_CONSOLE_TYPE_WARN:
				console.Warn("Warning!!")
			}
			return nil
		})

		pn.Console = console
	}

	pn.Console.Log(
		pn.Console.EOL(),
		pn.Console.Text(
			banner,
			utils.PN_CONSOLE_COLOR_GREEN,
			utils.PN_CONSOLE_COLOR_BLACK,
			utils.PN_CONSOLE_STYLE_BOLD),
	)

	// load env_module
	if err := pn.Load(&PnDotEnv{}); err != nil {

		pn.Console.Error(err)
	}

	if pn.Echo == nil {

		pn.Echo = echo.New()
		//pn.Echo.Pre(middleware.Logger())
		//pn.Echo.Use(middleware.Recover())
		pn.Echo.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(ctx echo.Context) error {

				pn.Console.Main(func() error {

					request := ctx.Request()
					URL := net2.PnGetURLFromRequest(request)
					pn.Console.Log(URL) // Logger
					return nil
				})

				return next(ctx)
			}
		})
		pn.Echo.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(ctx echo.Context) error {

				response := ctx.Response()
				response.Header().Set("Server", "PapayaNet "+pn.Version.Stringify())

				return next(ctx)
			}
		})
	}

	if pn.DB == nil {

		client := &mysql.DBClient{}
		myClient, err := client.Init(db.DB_INIT_LOAD_ENV)
		if err != nil {

			pn.Console.Error(err)
		}
		if myClient != nil {

			pn.DB = myClient
		}
	}

	pn.Echo.HideBanner = true
	if echoLog, ok := pn.Echo.Logger.(*gLog.Logger); ok {

		echoLog.SetHeader("[${time_rfc3339_nano}] [${level}]")
	}

	pn.Console.Error("testing ...")

	return nil
}

func (pn *PapayaNet) Load(environ PnEnvImpl) error {

	return environ.Load()
}

func (pn *PapayaNet) Serve(addr string, port uint16) error {

	if port < 1024 {

		return errors.New("params `port` prevent to use! that mean port use for core systems")
	}

	sPort := strconv.Itoa(int(port))

	pHttp := &http.Server{Addr: addr + ":" + sPort}
	pHttp.Handler = pn.Echo
	pHttp.ErrorLog = pn.Echo.StdLogger

	if pHttp.TLSConfig == nil {

		Listener, err := net.Listen("tcp", pHttp.Addr)

		if err != nil {

			return err
		}

		pn.Echo.Listener = Listener

		pn.Console.Log("Server started on " + addr + ":" + sPort)
		return pHttp.Serve(pn.Echo.Listener)
	}

	TLSListener, err := tls.Listen("tcp", pHttp.Addr, pHttp.TLSConfig)

	if err != nil {

		return err
	}

	pn.Echo.TLSListener = TLSListener

	pn.Console.Log("Server started on " + addr + ":" + sPort)
	return pHttp.Serve(pn.Echo.TLSListener)
}
