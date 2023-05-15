package easy

import "gorm.io/gorm"

type Model struct {
  *gorm.Model
  ID string `gorm:"type:VARCHAR(32);primary" json:"id"`
}

func (Model) TableName() string {

  return "model"
}
