package models

import (
  "github.com/google/uuid"
  "gorm.io/gorm"
  "time"
)

type SessionModel struct {
  *gorm.Model
  ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary" json:"id"`
  Token     string    `gorm:"type:text;unique;not null" json:"token"`
  SecretKey string    `gorm:"type:text;unique;not null" json:"secret_key"`
  Expired   time.Time `gorm:"type:timestamp;not null" json:"expired"`
}

// set table name

func (SessionModel) TableName() string {

  return "sessions"
}
