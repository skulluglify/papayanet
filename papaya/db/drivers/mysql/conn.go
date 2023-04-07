package mysql

import (
	"PapayaNet/papaya/db"
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type DBClient struct {
	Host     string
	Port     uint16 // 65535
	User     string
	Pass     string
	Name     string
	Charset  string // utf8mb4
	TimeZone string // UTC
	//Secure   bool   // allow, disable, ...
}

type DBConfigImpl interface {
	Stringify() string
	IsUnixSock(path string) bool
	Init(flags int) (*gorm.DB, error)
}

func (client *DBClient) IsUnixSock() bool {

	if strings.HasPrefix(client.Host, "/") ||
		strings.Contains(client.Host, "/") ||
		strings.HasSuffix(client.Host, ".sock") {

		return true
	}

	return false
}

func (client *DBClient) Stringify() string {

	//SSLMode := "disable"
	//
	//if client.Secure {
	//
	//	SSLMode = "enable"
	//}

	params := url.Values{}
	params.Set("charset", client.Charset)
	//params.Set("sslmode", SSLMode)
	params.Set("parseTime", "true")
	if loc, err := time.LoadLocation(client.TimeZone); err == nil {

		params.Set("loc", loc.String())
	} else {

		params.Set("loc", "UTC")
	}

	userInfo := url.UserPassword(client.User, client.Pass)
	URL := url.URL{
		User: userInfo,
		Host: func() string {

			if client.Port > 1023 {

				return "tcp(" + client.Host + ":" + strconv.Itoa(int(client.Port)) + ")"

			} else {

				log.Fatal("var `port` prevent to use! that mean port use for core systems")
			}

			if //goland:noinspection ALL
			client.IsUnixSock() {

				return "unix(" + client.Host + ")"
			}

			return "tcp(" + client.Host + ":3306" + ")"
		}(),
		Scheme:   "mysql",
		Path:     client.Name,
		RawQuery: params.Encode(),
	}

	return URL.String()
}

func (client *DBClient) Init(flags int) (*gorm.DB, error) {
	if flags|db.DB_INIT_LOAD_ENV == 1 {

		port, err := strconv.ParseUint(os.Getenv("DB_PORT"), 10, 16)

		if err != nil {
			return nil, errors.New("cannot parse `port` as `uint16`")
		}

		client.Host = os.Getenv("DB_HOST")
		client.User = os.Getenv("DB_USER")
		client.Port = uint16(port)
		client.Pass = os.Getenv("DB_PASS")
		client.Charset = os.Getenv("DB_CHARSET")
		client.Name = os.Getenv("DB_NAME")
		client.TimeZone = os.Getenv("DB_TIMEZONE")

		//if sec, err := strconv.ParseBool(os.Getenv("DB_SECURE")); err == nil {
		//
		//	client.Secure = sec
		//}
	}

	if client.Host == "" {

		return nil, errors.New("var `host` not found")
	}

	dsn, _ := url.QueryUnescape(client.Stringify()[8:])

	myClient, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{})

	if err != nil {
		return nil, errors.New("cannot create `db_mysql_client`")
	}

	return myClient, nil
}

// func select all, select, insert, update, delete, delete soft, select UnScope
