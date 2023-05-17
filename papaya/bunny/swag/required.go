package swag

import "skfw/papaya/koala"

func SwagRequired(h string) (bool, string) {

  n := len(h)

  if koala.KStrHasPrefixChar(h, "?") {

    return false, h[1:]
  }

  if koala.KStrHasSuffixChar(h, "?") {

    return false, h[:n-1]
  }

  if koala.KStrHasPrefixChar(h, "!") {

    return true, h[1:]
  }

  if koala.KStrHasSuffixChar(h, "!") {

    return true, h[:n-1]
  }

  // default, required
  return true, h
}
