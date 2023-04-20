package swag

import (
  m "PapayaNet/papaya/koala/mapping"
  "PapayaNet/papaya/koala/pp"
)

//AuthToken
//request.validation

//description
//request.params
//request.headers
//request.body.application/json
//responses.200
//responses.200.body.application/json

type SwagExpect struct {
  AuthToken         bool
  RequestValidation bool
  Path              m.KMapImpl
}

func SwagExpectEvaluation(expect m.KMapImpl, tags []string) *SwagExpect {

  var requestBody m.KMapImpl

  reqVal := expect.Get("request.validation")
  checkRequestValidation := reqVal != nil

  authToken := m.KValueToBool(expect.Get("AuthToken")) // default value is false
  description := m.KValueToString(expect.Get("description"))
  requestValidation := pp.LBool(checkRequestValidation, m.KValueToBool(reqVal), true) // default value is true
  parameters := SwagParamsFormatter(expect.Get("request.params"))
  body := m.KMapCast(expect.Get("request.body"))
  res := m.KMapCast(expect.Get("responses"))

  requestBody = nil
  bodyContentSchemes := SwagContentSchemes(body)

  if len(bodyContentSchemes) > 0 {

    requestBody = bodyContentSchemes[0]
  }

  path := &m.KMap{

    "description": description,
    "tags":        tags,
    "responses":   SwagResponseSchemes(res),
  }

  // suitable for available values
  if len(parameters) > 0 {

    path.Put("parameters", parameters)
  }

  // suitable for available value
  if requestBody != nil {

    path.Put("requestBody", requestBody)
  }

  return &SwagExpect{
    AuthToken:         authToken,
    RequestValidation: requestValidation,
    Path:              path,
  }
}
