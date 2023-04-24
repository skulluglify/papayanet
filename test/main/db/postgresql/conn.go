package main

import (
  "PapayaNet/papaya/pigeon"
  "PapayaNet/papaya/pigeon/drivers/postgresql"
  "PapayaNet/papaya/pigeon/templates/basic/models"
  "log"
  "os"
)

func main() {

  os.Setenv("DB_HOST", "localhost")
  os.Setenv("DB_PORT", "5432")
  os.Setenv("DB_USERNAME", "user")
  os.Setenv("DB_PASSWORD", "1234")
  os.Setenv("DB_NAME", "leafy")
  os.Setenv("DB_CHARSET", "utf8")
  os.Setenv("DB_TIMEZONE", "UTC")
  os.Setenv("DB_SECURE", "false")
  os.Setenv("DB_UNIX_SOCK", "")

  conn, err := postgresql.DBConnectionNew(pigeon.InitLoadEnviron)
  defer func(conn *postgresql.DBConnection) {
    err := conn.Close()
    if err != nil {
      log.Fatal(err)
    }
  }(conn)

  if err != nil {

    log.Fatal(err)
  }

  if err := postgresql.PgEnableExtensionUUID(conn); err != nil {

    log.Fatal(err)
  }

  if err := postgresql.PgSetTimeZoneUTC(conn); err != nil {

    log.Fatal(err)
  }

  err = conn.DB.AutoMigrate(&models.UserModel{}, &models.SessionModel{})
  if err != nil {

    log.Fatal(err)
  }

}
