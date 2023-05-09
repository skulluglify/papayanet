package swag

import (
	m "skfw/papaya/koala/mapping"
	"skfw/papaya/koala/pp"
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
	RequestValidation bool
	Path              m.KMapImpl
}

func SwagExpectEvaluation(expect m.KMapImpl, tags []string) *SwagExpect {

	var requestBody m.KMapImpl

	// request validation expected
	// nil -> true
	// true -> true
	// false -> false
	reqVal := expect.Get("request.validation")
	checkRequestValidation := reqVal != nil

	description := m.KValueToString(expect.Get("description"))
	requestValidation := pp.Lbool(checkRequestValidation, m.KValueToBool(reqVal), true) // default value is true
	parameters := SwagParamsFormatter(expect.Get("request.params"))
	headers := SwagHeadersFormatter(expect.Get("request.headers"))
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

	// append headers in the parameters
	if len(headers) > 0 {

		for _, header := range headers {

			parameters = append(parameters, header)
		}
	}

	// suitable for available values
	if len(parameters) > 0 {

		path.Put("parameters", parameters)
	}

	// suitable for available request body
	if requestBody != nil {

		path.Put("requestBody", requestBody)
	}

	return &SwagExpect{
		RequestValidation: requestValidation,
		Path:              path,
	}
}
