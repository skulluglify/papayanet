package swag

import (
  "skfw/papaya/koala/kornet"
  "strings"
)

// compare key from enums

// do.main.x.register -> .x. -> array sample

type CompareSampleKey struct {
  tokens []string
}

type CompareSampleKeyImpl interface {
  Init(key string) error
  Check(key string) bool
}

func CompareSampleKeyNew(key string) (CompareSampleKeyImpl, error) {

  compareSampleKey := &CompareSampleKey{}
  if err := compareSampleKey.Init(key); err != nil {

    return nil, err
  }
  return compareSampleKey, nil
}

func (c *CompareSampleKey) Init(key string) error {

  c.tokens = strings.Split(key, ".")
  return nil
}

func (c *CompareSampleKey) Check(key string) bool {

  var token string
  var tokens []string
  var n int

  tokens = strings.Split(key, ".")
  n = len(c.tokens)

  if n != len(tokens) {

    return false
  }

  for i := 0; i < n; i++ {

    token = c.tokens[i]

    if token == "0" {

      if _, err := kornet.KSafeParsingNumber(tokens[i]); err != nil {

        return false
      }

      continue
    }

    if token != tokens[i] {

      return false
    }
  }

  return true
}
