package test

import (
  "PapayaNet/papaya"
  "testing"
)

func Test(test *testing.T) {

  pn := &papaya.PapayaNet{}
  err := pn.Init()
  if err != nil {
    return
  }

  test.Log("Initial Completed ...")
}
