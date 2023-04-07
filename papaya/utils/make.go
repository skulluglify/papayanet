package utils

import (
	"strconv"
	"strings"
)

type PnVersion struct {
	Major uint8
	Minor uint8
	Patch uint8
}

type PnVersionImpl interface {
	Stringify() string
}

func PnMakeVersion(major uint8, minor uint8, patch uint8) *PnVersion {

	return &PnVersion{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}

func (version *PnVersion) Stringify() string {

	return "v" + strings.Join([]string{
		strconv.Itoa(int(version.Major)),
		strconv.Itoa(int(version.Minor)),
		strconv.Itoa(int(version.Patch))},
		".")
}
