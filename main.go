package main

import (
	"PapayaNet/app"
	"PapayaNet/papaya"
	"fmt"
)

func main() {

	fmt.Println("Papaya Net v1.0 testing ...")

	pn := papaya.NetNew()
	pn.Init()

	if err := app.App(pn); err != nil {

		pn.Logger().Error(err)
	}

	if err := pn.Close(); err != nil {

		pn.Logger().Error(err)
	}

	pn.Logger().Log("Shutdown ...")
}
