package mysql

import (
  "database/sql"
  "errors"
  my "github.com/go-sql-driver/mysql"
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
  "math"
  "net/url"
  "skfw/papaya/koala/environ"
  "skfw/papaya/koala/kio"
  "skfw/papaya/koala/pp"
  "skfw/papaya/pigeon"
  "skfw/papaya/pigeon/drivers/common"
  "time"
)

// ref: https://github.com/go-sql-driver/mysql
// ref: https://dev.mysql.com/doc/refman/8.0/en/grant.html#grant-database-privileges

type DBConfig struct {
  *common.DBConfig

  // mysql var config
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
        Port: 3306,
      },
    },
  }

  _, err := conn.Init(flags)
  return conn, err
}

func (c *DBConnection) Init(flags int) (*gorm.DB, error) {

  var err error

  if pp.Qflag(flags, pigeon.InitLoadEnviron) {

    // common load
    environ.KEnvLoaderNew[*common.DBConfig]().Load(c.DBConfig.DBConfig)

    // current load
    environ.KEnvLoaderNew[*DBConfig]().Load(c.DBConfig)
  }

  if math.MaxUint16 < c.DBConfig.Port {

    panic("var `port` has been configured incorrectly")
  }

  DSN := c.DSN()

  c.DB, err = gorm.Open(mysql.Open(DSN), c.Config)

  if err != nil {

    return nil, err
  }

  // make it UTC mode
  if err = c.DB.Exec("SET TIME_ZONE = ?", "+00:00").Error; err != nil {

    return nil, err
  }

  return c.DB, nil
}

func (c *DBConnection) IsUnixSock() bool {

  return kio.KFileNew(c.UnixSock).IsSocket()
}

func (c *DBConnection) RawQuery() string {

  values := url.Values{}

  // mysql idk

  return values.Encode()
}

func (c *DBConnection) Database() (*sql.DB, error) {

  return c.DB.DB()
}

func (c *DBConnection) GORM() *gorm.DB {

  return c.DB
}

func (c *DBConnection) DSN() string {

  // read pass file as var `password`
  if c.DBConfig.Password == "" {

    f := kio.KFileNew(c.DBConfig.PassFile)

    if f.IsExist() {

      c.DBConfig.Password = f.Cat()
    }
  }

  config := my.Config{
    Net:       pp.Lstr(c.IsUnixSock(), "unix", "tcp"),
    Addr:      c.DBConfig.Host,
    User:      c.DBConfig.Username,
    Passwd:    c.DBConfig.Password,
    DBName:    c.DBConfig.Name,
    ParseTime: true,
    Loc:       time.UTC,
  }

  return config.FormatDSN()
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
