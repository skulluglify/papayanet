package models

import (
  "time"

  "gorm.io/gorm"
)

// ID

// uuid.UUID generate 32 chars, and convert into 16 bytes

// Client-IP

// In general, IPv4 addresses are 32 bits long and are typically represented as a string of four decimal numbers,
// separated by periods (e.g. "192.0.2.1"). IPv6 addresses, on the other hand,
// are 128 bits long and are typically represented as a string of hexadecimal digits,
//  separated by colons (e.g. "2001:0db8:85a3:0000:0000:8a2e:0370:7334").

// The longest IPv4 address possible is 255.255.255.255, which is 15 characters long.
//  The longest IPv6 address possible is 8 groups of 4 hexadecimal digits,
//  separated by colons, which is 39 characters long.

// User-Agent

// The user agent string for the Tor Browser, which can be up to 1024 characters long,
// due to the addition of extra information about the user's privacy settings.
// The user agent string for Google Chrome on iOS, which can be up to 8192 characters long,
// due to the addition of extra debugging information.

// MySQL, SQLite not support UUID datatype
// PostgreSQL not support BINARY, or VARBINARY datatype

// set into HEX

// Token -> HashToken
// SecretKey -> 32 bytes into 44 chars (base64)

type SessionModel struct {
  *gorm.Model
  ID            string    `gorm:"type:VARCHAR(32);primary" json:"id"`
  UserID        string    `gorm:"type:VARCHAR(32);not null" json:"user_id"`
  ClientIP      string    `gorm:"type:VARCHAR(40);not null" json:"client_ip"`
  UserAgent     string    `gorm:"type:TEXT;not null" json:"user_agent"`
  Token         string    `gorm:"type:VARCHAR(64);unique;not null" json:"token"`
  SecretKey     string    `gorm:"type:VARCHAR(44);unique;not null" json:"secret_key"`
  Expired       time.Time `gorm:"type:TIMESTAMP;not null" json:"expired"`
  LastActivated time.Time `gorm:"type:TIMESTAMP;default:CURRENT_TIMESTAMP" json:"last_activated"`
}

// set table name

func (SessionModel) TableName() string {

  return "sessions"
}
