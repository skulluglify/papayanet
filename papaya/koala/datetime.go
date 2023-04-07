package koala

import (
  "strconv"
  "time"
)

type KDateTime struct {
  Year        int
  Month       int
  Day         int
  Hour        int
  Minute      int
  Second      int
  Microsecond int
  Offset      int
  Zone        string
}

type KDateTimeImpl interface {
  UTC() *KDateTime
  String() string
  Ez() string
}

func KDateTimeNew() *KDateTime {

  // ISO 8601 format
  local := time.Now()
  zone, offset := local.Zone()

  return &KDateTime{
    Year:        local.Year(),
    Month:       int(local.Month()),
    Day:         local.Day(),
    Hour:        local.Hour(),
    Minute:      local.Minute(),
    Second:      local.Second(),
    Microsecond: local.Nanosecond() / 1e6,
    Offset:      offset,
    Zone:        zone,
  }
}

func (t *KDateTime) UTC() *KDateTime {

  utc := time.Now().UTC()
  zone, offset := utc.Zone()

  return &KDateTime{
    Year:        utc.Year(),
    Month:       int(utc.Month()),
    Day:         utc.Day(),
    Hour:        utc.Hour(),
    Minute:      utc.Minute(),
    Second:      utc.Second(),
    Microsecond: utc.Nanosecond() / 1e6,
    Offset:      offset,
    Zone:        zone,
  }
}

func (t *KDateTime) String() string {

  Y := KStrZeroFill(strconv.Itoa(t.Year), 4)
  M := KStrZeroFill(strconv.Itoa(t.Month), 2)
  D := KStrZeroFill(strconv.Itoa(t.Day), 2)
  H := KStrZeroFill(strconv.Itoa(t.Hour), 2)
  m := KStrZeroFill(strconv.Itoa(t.Minute), 2)
  s := KStrZeroFill(strconv.Itoa(t.Second), 2)
  f := KStrZeroFill(strconv.Itoa(t.Microsecond), 3)

  return Y + "-" + M + "-" + D + "T" + H + ":" + m + ":" + s + "." + f + "Z"
}

func (t *KDateTime) Ez() string {

  Y := KStrZeroFill(strconv.Itoa(t.Year), 4)
  M := KStrZeroFill(strconv.Itoa(t.Month), 2)
  D := KStrZeroFill(strconv.Itoa(t.Day), 2)
  H := KStrZeroFill(strconv.Itoa(t.Hour), 2)
  m := KStrZeroFill(strconv.Itoa(t.Minute), 2)
  s := KStrZeroFill(strconv.Itoa(t.Second), 2)

  return Y + "/" + M + "/" + D + " " + H + ":" + m + ":" + s + " " + t.Zone
}
