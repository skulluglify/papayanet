package main

import (
  "PapayaNet/papaya/db"
  "PapayaNet/papaya/db/drivers/postgresql"
  "log"
)

func main() {

  conn, err := postgresql.DBConnectionNew(db.InitLoadEnviron)
  if err == nil {
    err = conn.Close()
  }
  log.Fatal(err)
}
