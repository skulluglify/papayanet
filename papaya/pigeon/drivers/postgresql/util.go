package postgresql

import (
  "database/sql"
  "skfw/papaya/pigeon/drivers/common"
  "strings"
)

func Query(conn common.DBConnectionImpl, query string, args ...any) (*sql.Rows, error) {

  // SHORT NAME function query FROM db.DB()

  db, err := conn.Database()

  if err != nil {

    return nil, err
  }

  return db.Query("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
}

func IsPostgreSQL(conn common.DBConnectionImpl) bool {

  return strings.HasPrefix(conn.DSN(), "postgres")
}

func PgEnableExtensionUUID(conn common.DBConnectionImpl) error {

  // FIX PROBLEM USE function uuid_generate_v4()

  // check scheme is postgres
  if IsPostgreSQL(conn) {

    if _, err := Query(conn, "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"); err != nil {

      return err
    }
  }

  return nil
}

func PgSetTimeZoneUTC(conn common.DBConnectionImpl) error {

  // GENERIC PURPOSE SET TIMEZONE INTO UTC MODE

  // check scheme is postgres
  if IsPostgreSQL(conn) {

    if _, err := Query(conn, "SET TIME ZONE \"UTC\";"); err != nil {

      return err
    }
  }

  return nil
}
