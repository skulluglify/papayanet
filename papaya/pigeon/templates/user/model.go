package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	*gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
	Password string    `json:"password"`
	Address  string    `json:"address"`
	Country  string    `json:"country"`
	City     string    `json:"city"`
	DOB      time.Time `json:"dob"`
}
