package main

import (
  "fmt"
  "os"
  "skfw/papaya/pigeon/drivers/mysql"
  "skfw/papaya/pigeon/templates/basicAuth/util"
)

func main() {

  os.Setenv("DB_NAME", "main")
  os.Setenv("DB_USERNAME", "user")
  os.Setenv("DB_PASSWORD", "User@1234")

  conn, _ := mysql.DBConnectionNew(0)
  fmt.Println(conn.DSN())
  fmt.Println(util.CreateSecretKey())
}
