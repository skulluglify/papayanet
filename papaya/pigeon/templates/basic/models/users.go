package models

import (
  "github.com/google/uuid"
  "gorm.io/gorm"
  "time"
)

// ref: https://gorm.io/docs/models.html#Fields-Tags

// Username
// The maximum length of a username for a database varies depending on the specific database system being used
// Some databases may have a maximum username length of 16 characters, while others may allow usernames to be up to 32 characters long or even longer.

// A username is a unique identifier used to log into an account or access a system.
// It is typically shorter than a person's full name and may include letters, numbers, and special characters.
// A long name, on the other hand, typically refers to a person's full name, which may include their first name, middle name(s), and last name.
// A long name is generally longer than a username and may include spaces and punctuation.
// The main difference between a username and a long name is their purpose:
// - a username is used for identification and authentication, while a long name is used to identify a person.

// EMail
// The maximum length of an email address is defined by the RFC 5321 standard as 254 characters.
// This includes the local part (the part before the '@' symbol) and the domain part (the part after the '@' symbol).

// Password
// The maximum length of a password for a database varies
// Depending on the specific database system being used.
// Some databases have a maximum password length of 32 characters,
// while others may allow passwords to be up to 128 characters long or even longer.
// It is generally recommended to use a password that is at least 8 characters long
// Includes a combination of upper and lower case letters, numbers, and special characters to increase security.

// Gender
// Male or Female, not Other anymore

// Phone Number
// The International Telecommunication Union (ITU)
// Recommends that phone numbers have a maximum length of 15 digits,
// Including the country code.
// However, some countries have phone numbers that are shorter or longer than this recommendation.

// Country codes are typically two or three letters long.
// The International Organization for Standardization (ISO) has established a standard for country codes known as ISO 3166.
// According to this standard, country codes can be either two letters (ISO 3166-1 alpha-2) or three letters (ISO 3166-1 alpha-3) long.
// So the maximum length of a country code according to the ISO standard is three letters.

// City Name
// ref: https://largest.org/geography/city-names/

// Postal Code
// ref: https://www.grcdi.nl/pidm/postal%20code.html

type UserModel struct {
  *gorm.Model
  ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary" json:"id"`
  Name        string    `gorm:"type:varchar(52);" json:"name"`
  Username    string    `gorm:"type:varchar(16);unique;not null" json:"username"`
  Email       string    `gorm:"type:varchar(254);unique;not null" json:"email"`
  Password    string    `gorm:"type:varchar(128);not null" json:"password"`
  Gender      string    `gorm:"type:char(1)" json:"gender"`
  Phone       string    `gorm:"type:varchar(24)" json:"phone"`
  DOB         time.Time `gorm:"type:timestamp" json:"dob"`
  Address     string    `gorm:"type:varchar(128)" json:"address"`
  CountryCode string    `gorm:"type:varchar(4)" json:"country_code"`
  City        string    `gorm:"type:varchar(64)" json:"city"`
  PostalCode  string    `gorm:"type:varchar(10)" json:"postal_code"`
  Admin       bool      `gorm:"type:boolean" gorm:"default:false" json:"admin"`
}

// set table name

func (UserModel) TableName() string {

  return "users"
}