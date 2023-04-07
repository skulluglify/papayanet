package main

import (
	"PapayaNet/app"
	"PapayaNet/papaya"
	"fmt"
	"os"
)

func main() {

	fmt.Println("Papaya Net v1.0 testing ...")

	pn := papaya.PapayaNet{}
	if err := pn.Init(); err != nil {

		panic(err)
	}

	if err := app.App(&pn); err != nil {

		pn.Console.Error(err)
		os.Exit(1)
	}

	if err := pn.Serve("127.0.0.1", 8000); err != nil {

		pn.Console.Error(err)
	}
}
