package easy

import (
  "gorm.io/gorm"
  "time"
)

type Model struct {
  ID        string         `gorm:"type:VARCHAR(32);primary" json:"id"`
  CreatedAt time.Time      `json:"created_at"`
  UpdatedAt time.Time      `json:"updated_at"`
  DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Model) TableName() string {

  return "model"
}
