package user

import (
  "github.com/google/uuid"
  "gorm.io/gorm"
)

type User struct {
  *gorm.Model
  ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
  Name     string
  Email    string
  Phone    string
  Password string
  Address  string
  Country  string
  City     string
}
