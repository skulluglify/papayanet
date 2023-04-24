package postgresql

import "database/sql"

func Query(conn DBConnectionImpl, query string, args ...any) (*sql.Rows, error) {

  // SHORT NAME function query FROM db.DB()

  db, err := conn.Database()

  if err != nil {

    return nil, err
  }

  return db.Query("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
}

func PgEnableExtensionUUID(conn DBConnectionImpl) error {

  // FIX PROBLEM USE function uuid_generate_v4()

  if _, err := Query(conn, "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"); err != nil {

    return err
  }

  return nil
}

func PgSetTimeZoneUTC(conn DBConnectionImpl) error {

  // GENERIC PURPOSE SET TIMEZONE INTO UTC MODE

  if _, err := Query(conn, "SET TIME ZONE \"UTC\";"); err != nil {

    return err
  }

  return nil
}
