package main

import (
  "PapayaNet/papaya/db/drivers/postgresql"
  "log"
)

func main() {

	conn, err := postgresql.DBConnectionNew(pigeon.InitLoadEnviron)
	if err == nil {
		err = conn.Close()
	}
	log.Fatal(err)
}
