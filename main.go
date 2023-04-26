package main

import (
	"fmt"
	"skfw/app"
	"skfw/papaya"
)

func main() {

	fmt.Println("Papaya Net v1.0 testing ...")

	pn := papaya.NetNew()

	if err := app.App(pn); err != nil {

		pn.Logger().Error(err)
	}

	if err := pn.Close(); err != nil {

		pn.Logger().Error(err)
	}

	pn.Logger().Log("Shutdown ...")
}
