package utils

import (
	"strconv"
	"time"
)

type PnTime struct {
	Year        int
	Month       int
	Day         int
	Hour        int
	Minute      int
	Second      int
	Microsecond int
}

type PnTimeImpl interface {
	Stringify() string
	Simple() string
}

func PnTimeNow() *PnTime {

	// ISO 8601 format
	t := time.Now().UTC()
	return &PnTime{
		Year:        t.Year(),
		Month:       int(t.Month()),
		Day:         t.Day(),
		Hour:        t.Hour(),
		Minute:      t.Minute(),
		Second:      t.Second(),
		Microsecond: t.Nanosecond() / 1e6,
	}
}

func (t *PnTime) Stringify() string {

	Y := PnStrZeroFill(strconv.Itoa(t.Year), 4)
	M := PnStrZeroFill(strconv.Itoa(int(t.Month)), 2)
	D := PnStrZeroFill(strconv.Itoa(t.Day), 2)
	H := PnStrZeroFill(strconv.Itoa(t.Hour), 2)
	m := PnStrZeroFill(strconv.Itoa(t.Minute), 2)
	s := PnStrZeroFill(strconv.Itoa(t.Second), 2)
	f := PnStrZeroFill(strconv.Itoa(t.Microsecond), 3)

	return Y + "-" + M + "-" + D + "T" + H + ":" + m + ":" + s + "." + f + "Z"
}

func (t *PnTime) Simple() string {

	Y := PnStrZeroFill(strconv.Itoa(t.Year), 4)
	M := PnStrZeroFill(strconv.Itoa(int(t.Month)), 2)
	D := PnStrZeroFill(strconv.Itoa(t.Day), 2)
	H := PnStrZeroFill(strconv.Itoa(t.Hour), 2)
	m := PnStrZeroFill(strconv.Itoa(t.Minute), 2)
	s := PnStrZeroFill(strconv.Itoa(t.Second), 2)

	return Y + "/" + M + "/" + D + " " + H + ":" + m + ":" + s
}
