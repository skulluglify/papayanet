package postgresql

import (
	"PapayaNet/papaya/koala/environ"
	"PapayaNet/papaya/koala/kio"
	"PapayaNet/papaya/koala/pp"
	"PapayaNet/papaya/pigeon"
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"math"
	"net/url"
	"strconv"
	"strings"
)

// Ref: https://www.postgresql.org/docs/current/libpq-connect.html
// Ref: https://www.postgresql.org/docs/15/runtime-config-client.html
// Ref: https://www.prisma.io/docs/concepts/database-connectors/postgresql

// TODO: postgresql not implemented connection by socket

type DBConfig struct {

	// general purpose
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT"`
	UnixSock string `env:"DB_UNIX_SOCK"`
	Username string `env:"DB_USERNAME"`
	Password string `env:"DB_PASSWORD"`
	PassFile string `env:"DB_PASSWORD_FILE"`
	TimeZone string `env:"DB_TIMEZONE"`
	Charset  string `env:"DB_CHARSET"`
	Secure   bool   `env:"DB_SECURE"`
	Name     string `env:"DB_NAME"`

	// postgresql only
	PgConnectTimeout  int    `env:"PG_CONNECT_TIMEOUT"`
	PgSSLMode         string `env:"PG_SSL_MODE"`
	PgSSLCert         string `env:"PG_SSL_CERT"`
	PgSSLPassword     string `env:"PG_SSL_PASSWORD"`
	PgApplicationName string `env:"PG_APPLICATION_NAME"`
	PgOptions         string `env:"PG_OPTIONS"`
	PgClientEncoding  string `env:"PG_CLIENT_ENCODING"`
	PgTimeZone        string `env:"PG_TIMEZONE"`
}

type DBConnection struct {
	*gorm.Config
	*gorm.DB
	*DBConfig
}

type DBConfigImpl interface {
	Init(flags int) (*gorm.DB, error)
	IsUnixSock() bool
	RawQuery() string
	String() string
	Close() error
}

func DBConnectionNew(flags int) (*DBConnection, error) {

	conn := &DBConnection{
		Config: &gorm.Config{},
		DBConfig: &DBConfig{
			Port: 5432,
		},
	}

	_, err := conn.Init(flags)
	return conn, err
}

func (c *DBConnection) Init(flags int) (*gorm.DB, error) {

	if pp.KValidFlag(flags, pigeon.InitLoadEnviron) {

		envLoader := environ.KEnvLoaderNew[*DBConfig]()
		envLoader.Load(c.DBConfig)
	}

	if math.MaxUint16 < c.Port {

		panic("var `port` has been configured incorrectly")
	}

	DSN := c.String()

	DB, err := gorm.Open(postgres.Open(DSN), c)

	if err != nil {

		return nil, err
	}

	c.DB = DB
	return DB, nil
}

func (c *DBConnection) IsUnixSock() bool {

	// that may problem, bcs postgresql not use literally sock file to connect them
	return kio.KFileNew(c.UnixSock).IsSocket()
}

func (c *DBConnection) RawQuery() string {

	values := url.Values{}
	connectTimeout := 5
	sslmode := "prefer" // disable|prefer|require
	sslcert := ""
	sslpassword := ""
	applicationName := ""
	options := ""
	clientEncoding := "utf8"
	timeZone := "UTC"

	// general purpose
	if c.Secure {

		sslmode = "require"
	}

	if c.Charset != "" {

		clientEncoding = c.Charset
	}

	if c.TimeZone != "" {

		timeZone = c.TimeZone
	}

	// postgresql only

	if c.PgConnectTimeout != 0 {
		connectTimeout = c.PgConnectTimeout
	}

	if c.PgSSLMode != "" {
		sslmode = c.PgSSLMode
	}

	if c.PgSSLCert != "" {
		sslcert = c.PgSSLCert
	}

	if c.PgSSLPassword != "" {
		sslpassword = c.PgSSLPassword
	}

	if c.PgApplicationName != "" {
		applicationName = c.PgApplicationName
	}

	if c.PgOptions != "" {
		options = c.PgOptions
	}

	if c.PgClientEncoding != "" {
		clientEncoding = c.PgClientEncoding
	}

	if c.PgTimeZone != "" {
		timeZone = c.PgTimeZone
	}

	values.Add("connect_timeout", strconv.Itoa(connectTimeout))
	values.Add("sslmode", sslmode)
	values.Add("sslcert", sslcert)
	values.Add("sslpassword", sslpassword)
	values.Add("application_name", applicationName)
	values.Add("options", options)
	values.Add("client_encoding", clientEncoding)
	values.Add("TimeZone", timeZone)

	return values.Encode()
}

func (c *DBConnection) String() string {

	//// trim all strings
	//c.Host = strings.Trim(c.Host, " ")
	//c.UnixSock = strings.Trim(c.UnixSock, " ")
	//c.Username = strings.Trim(c.Username, " ")
	//c.Password = strings.Trim(c.Password, " ")
	//c.PassFile = strings.Trim(c.PassFile, " ")
	//c.TimeZone = strings.Trim(c.TimeZone, " ")
	//c.Charset = strings.Trim(c.Charset, " ")
	//c.Name = strings.Trim(c.Name, " ")

	// set default value
	if c.Host != "" {

		c.Host = "localhost"
	}

	// replace var `port` with var `host`
	if c.Port != 0 {

		if !strings.Contains(c.Host, ":") {

			c.Host += ":" + strconv.Itoa(c.Port)
		}
	}

	// read pass file as var `password`
	if c.Password == "" {

		f := kio.KFileNew(c.PassFile)

		if f.IsExist() {

			c.Password = f.Cat()
		}
	}

	DSN := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.Username, c.Password),
		Host:     c.Host,
		Path:     c.Name,
		RawQuery: c.RawQuery(),
	}

	return DSN.String()
}

func (c *DBConnection) Close() error {

	if c.DB != nil {

		DB, err := c.DB.DB()
		if err != nil {

			return err
		}
		err = DB.Close()
		c.DB = nil
		return err
	}

	return errors.New("database has not been initialized")
}
