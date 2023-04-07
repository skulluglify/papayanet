package postgresql

import (
  "PapayaNet/papaya/db"
  "PapayaNet/papaya/koala/environ"
  "PapayaNet/papaya/koala/kio"
  "PapayaNet/papaya/koala/pp"
  "errors"
  "gorm.io/driver/postgres"
  "gorm.io/gorm"
  "net/url"
  "strconv"
)

type DBConfig struct {
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
}

type DBConnection struct {
  *gorm.Config
  *gorm.DB
  *DBConfig
}

type DBConfigImpl interface {
  Init(flags int) (*gorm.DB, error)
  IsUnixSock() bool
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

  if pp.KValidFlag(flags, db.InitLoadEnviron) {

    envLoader := environ.KEnvLoaderNew[*DBConfig]()
    envLoader.Load(c.DBConfig)
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

  return kio.KFileNew(c.UnixSock).IsExists()
}

func (c *DBConnection) String() string {

  var res string

  if kio.KFileNew(c.UnixSock).IsExists() {

    res += "host="
    res += c.UnixSock
    res += " "

  } else {

    res += "host="
    res += pp.KCOStr(c.Host, "127.0.0.1")
    res += " "

    res += "port="
    res += pp.KCOStr(strconv.Itoa(c.Port), "5432")
    res += " "
  }

  res += "user="
  res += pp.KCOStr(c.Username, "postgres")
  res += " "

  if c.Password == "" {

    f := kio.KFileNew(c.PassFile)

    if f.IsExists() {

      c.Password = f.Cat()
    }
  }

  res += "password="
  res += url.QueryEscape(c.Password)
  res += " "

  res += "dbname="
  res += pp.KCOStr(c.Name, "postgres")
  res += " "

  res += "client_encoding="
  res += pp.KCOStr(c.Charset, "utf8")
  res += " "

  res += "sslmode="
  res += pp.KISStr(c.Secure, "enable", "disable")
  res += " "

  res += "TimeZone="
  res += pp.KCOStr(c.TimeZone, "UTC")
  res += " "

  return res
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
