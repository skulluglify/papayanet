package utils

import (
	"PapayaNet/papaya/panda"
)

func PnStrZeroFill(text string, s int) string {

	var zeros string

	k := panda.Min(len(text), s)
	z := s - k

	for i := 0; i < z; i++ {

		zeros += "0"
	}

	return zeros + text[:k]
}

func PnStrPadStart(text string, s int) string {

	var pads string

	k := panda.Min(len(text), s)
	z := s - k

	for i := 0; i < z; i++ {

		pads += " "
	}

	return pads + text[:k]
}

func PnStrPadEnd(text string, s int) string {

	var pads string

	k := panda.Min(len(text), s)
	z := s - k

	for i := 0; i < z; i++ {

		pads += " "
	}

	return text[:k] + pads
}
