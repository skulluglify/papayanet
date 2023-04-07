package koala

import (
	"PapayaNet/papaya/koala"
	"testing"
)

func TestKMap(test *testing.T) {

	data := koala.KMap{
		"a": 12,
		"b": false,
		"c": &koala.KMap{
			"e": "completed",
		},
	}

	if koala.KMapGetValue[int]("a", data) == nil {

		test.Error("Getting data.a failed ...")
	}
}
