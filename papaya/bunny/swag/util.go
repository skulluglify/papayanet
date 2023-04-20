package swag

import (
  "PapayaNet/papaya/koala"
  m "PapayaNet/papaya/koala/mapping"
  "strings"
)

// normalize path

func SwagPathNorm(p string) string {

  res := make([]string, 0)
  tokens := strings.Split(p, "/")

  for _, token := range tokens {

    k := len(token)

    if k > 0 {

      if koala.KStrHasPrefixChar(token, ":") {

        res = append(res, "{"+token[1:]+"}")
        continue
      }
    }

    res = append(res, token)
  }

  return strings.Join(res, "/")
}

// short function

func OkJSON(data any, descriptions ...string) m.KMapImpl {

  return &m.KMap{
    "200": &m.KMap{
      "body": JSON(data, descriptions...),
    },
  }
}

// Created

func CreatedJSON(data any, descriptions ...string) m.KMapImpl {

  return &m.KMap{
    "201": &m.KMap{
      "body": JSON(data, descriptions...),
    },
  }
}

// passing description by variadic arguments

func JSON(data any, descriptions ...string) m.KMapImpl {

  description := "Ok"

  if len(descriptions) > 0 {

    description = descriptions[0]
  }

  return &m.KMap{
    "application/json": &m.KMap{
      "description": description,
      "schema":      data,
    },
  }
}
