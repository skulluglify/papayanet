package test

import (
  "skfw/papaya"
  "testing"
)

func Test(test *testing.T) {

  pn := &papaya.Net{}
  pn.Init()

  test.Log("Initial Completed ...")
}
