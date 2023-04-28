package postgresql

import (
  "database/sql"
  "errors"
  "gorm.io/driver/postgres"
  "gorm.io/gorm"
  "math"
  "net/url"
  "skfw/papaya/koala/environ"
  "skfw/papaya/koala/kio"
  "skfw/papaya/koala/pp"
  "skfw/papaya/pigeon"
  "skfw/papaya/pigeon/drivers/common"
  "strconv"
)

// Ref: https://www.postgresql.org/docs/current/libpq-connect.html
// Ref: https://www.postgresql.org/docs/15/runtime-config-client.html
// Ref: https://www.prisma.io/docs/concepts/database-connectors/postgresql

// TODO: postgresql not implemented connection by socket

type DBConfig struct {
  *common.DBConfig

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
  *gorm.DB
  *gorm.Config
  *DBConfig
}

func DBConnectionNew(flags int) (common.DBConnectionImpl, error) {

  conn := &DBConnection{
    Config: &gorm.Config{},
    DBConfig: &DBConfig{
      DBConfig: &common.DBConfig{
        Port: 5432,
      },
    },
  }

  _, err := conn.Init(flags)
  return conn, err
}

func (c *DBConnection) Init(flags int) (*gorm.DB, error) {

  if pp.QFlag(flags, pigeon.InitLoadEnviron) {

    // common load
    environ.KEnvLoaderNew[*common.DBConfig]().Load(c.DBConfig.DBConfig)

    // current load
    environ.KEnvLoaderNew[*DBConfig]().Load(c.DBConfig)
  }

  if math.MaxUint16 < c.DBConfig.Port {

    panic("var `port` has been configured incorrectly")
  }

  DSN := c.DSN()

  DB, err := gorm.Open(postgres.Open(DSN), c.Config)

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
  if c.DBConfig.Secure {

    sslmode = "require"
  }

  if c.DBConfig.Charset != "" {

    clientEncoding = c.DBConfig.Charset
  }

  if c.DBConfig.TimeZone != "" {

    timeZone = c.DBConfig.TimeZone
  }

  // postgresql only

  if c.DBConfig.PgConnectTimeout != 0 {
    connectTimeout = c.DBConfig.PgConnectTimeout
  }

  if c.DBConfig.PgSSLMode != "" {
    sslmode = c.DBConfig.PgSSLMode
  }

  if c.DBConfig.PgSSLCert != "" {
    sslcert = c.DBConfig.PgSSLCert
  }

  if c.DBConfig.PgSSLPassword != "" {
    sslpassword = c.DBConfig.PgSSLPassword
  }

  if c.DBConfig.PgApplicationName != "" {
    applicationName = c.DBConfig.PgApplicationName
  }

  if c.DBConfig.PgOptions != "" {
    options = c.DBConfig.PgOptions
  }

  if c.DBConfig.PgClientEncoding != "" {
    clientEncoding = c.DBConfig.PgClientEncoding
  }

  if c.DBConfig.PgTimeZone != "" {
    timeZone = c.DBConfig.PgTimeZone
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

func (c *DBConnection) Database() (*sql.DB, error) {

  return c.DB.DB()
}

func (c *DBConnection) GORM() *gorm.DB {

  return c.DB
}

func (c *DBConnection) DSN() string {

  return c.DBConfig.DSN("postgres", c.RawQuery())
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
