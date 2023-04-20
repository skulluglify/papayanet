package session

import (
  "github.com/google/uuid"
  "gorm.io/gorm"
  "time"
)

type User struct {
  *gorm.Model
  ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
  Token   string    `json:"token"`
  Expired time.Time `json:"expired"`
}
