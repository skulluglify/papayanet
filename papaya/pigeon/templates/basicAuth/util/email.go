package util

import (
	"errors"
	"net/mail"
	"strings"
)

var MailIsEmptyString = errors.New("mail address is empty string")
var InvalidEmailAddress = errors.New("invalid e-mail address")

type Email struct {
	address string
	tlds    []string
}

type EmailImpl interface {
	Init(address string) error
	Verify() (bool, error)
	Value() string
}

func EmailNew(address string) (EmailImpl, error) {

	email := &Email{}
	if err := email.Init(address); err != nil {

		return nil, err
	}
	return email, nil
}

func (m *Email) Init(address string) error {

	address = strings.Trim(address, " ")

	m.tlds = GetTLDs()

	if address != "" {

		m.address = address
		return nil
	}

	return MailIsEmptyString
}

func (m *Email) Verify() (bool, error) {

	if !TLDChecker(m.tlds, m.address) {

		return false, InvalidEmailAddress
	}

	if _, err := mail.ParseAddress(m.address); err != nil {

		return false, InvalidEmailAddress
	}

	return true, nil
}

func (m *Email) Value() string {

	return m.address
}
