package repository

import (
  "errors"
  "strings"
  "unicode"
)

// exceptions
var PasswordLongerThanSeventyTwoBytes = errors.New("password doesn't accept longer than 72 bytes")
var PasswordIsEmptyString = errors.New("password is empty string")
var PasswordIsTooShort = errors.New("password is too short")
var PasswordDoesNotContainSpecialCharacter = errors.New("password does not contain a special character")
var PasswordDoesNotContainNumberCharacter = errors.New("password does not contain a number character")
var PasswordDoestNotContainUpperCharacter = errors.New("password does not contain a upper character")
var PasswordDoesNotContainLowerCharacter = errors.New("password does not contain a lower character")

type Password struct {
  data    string
  special bool
  number  bool
  upper   bool
  lower   bool
  size    int
}

type PasswordImpl interface {
  Init(password string) error
  Verify(min int, special bool, num bool, upper bool, lower bool) (bool, error)
  Value() string
}

func PasswordNew(password string) (PasswordImpl, error) {

  pass := &Password{}
  if err := pass.Init(password); err != nil {

    return nil, err
  }
  return pass, nil
}

func (p *Password) Init(password string) error {

  password = strings.Trim(password, " ") // trim

  // unicode problem, convert into bytes
  if len([]byte(password)) > 72 {

    // bcrypt, request
    return PasswordLongerThanSeventyTwoBytes
  }

  if password != "" {
    p.data = password
    return nil
  }

  return PasswordIsEmptyString
}

func (p *Password) Verify(min int, special bool, num bool, upper bool, lower bool) (bool, error) {

  p.number = false
  p.special = false
  p.upper = false
  p.lower = false
  p.size = 0

  for _, c := range p.data {

    switch {
    case unicode.IsPunct(c) || unicode.IsSymbol(c):
      p.special = true

    case unicode.IsNumber(c):
      p.number = true

    case unicode.IsUpper(c):
      p.upper = true

    case unicode.IsLower(c):
      p.lower = true

    }

    p.size++
  }

  switch {
  case p.size < min:
    return false, PasswordIsTooShort

  case special && !p.special:
    return false, PasswordDoesNotContainSpecialCharacter

  case num && !p.number:
    return false, PasswordDoesNotContainNumberCharacter

  case upper && !p.upper:
    return false, PasswordDoestNotContainUpperCharacter

  case lower && !p.lower:
    return false, PasswordDoesNotContainLowerCharacter

  }

  return true, nil
}

func (p *Password) Value() string {

  return p.data
}
