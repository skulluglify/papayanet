package common

import (
  "database/sql"
  "gorm.io/gorm"
  "net/url"
  "skfw/papaya/koala/kio"
  "strconv"
  "strings"
)

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
}

type DBConfigImpl interface {
  Init() error
  DSN(scheme string, rawQuery string) string
}

type DBConnectionImpl interface {
  Init(flags int) (*gorm.DB, error)
  IsUnixSock() bool
  Database() (*sql.DB, error)
  GORM() *gorm.DB
  DSN() string
  RawQuery() string
  Close() error
}

func DBConfigNew(username string, password string, host string, port int, name string) (DBConfigImpl, error) {

  config := &DBConfig{
    Host:     host,
    Port:     port,
    Username: username,
    Password: password,
    Name:     name,
  }
  if err := config.Init(); err != nil {

    return nil, err
  }
  return config, nil
}

func (c *DBConfig) Init() error {

  return nil
}

func (c *DBConfig) DSN(scheme string, rawQuery string) string {

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
    Scheme:   scheme,
    User:     url.UserPassword(c.Username, c.Password),
    Host:     c.Host,
    Path:     c.Name,
    RawQuery: rawQuery,
  }

  return DSN.String()
}
