package test

import (
	"PapayaNet/papaya"
	"testing"
)

func Test(test *testing.T) {

	pn := &papaya.PapayaNet{}
	pn.Init()

	test.Log("Initial Completed ...")
}
