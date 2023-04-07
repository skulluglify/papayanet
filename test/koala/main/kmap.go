package main

import (
	"PapayaNet/papaya/koala"
	"fmt"
)

func main() {

	fmt.Println("KMap testing ...")

	data := koala.KMap{
		"version": "v1.0.0",
		"root": &koala.KMap{
			"systems": []string{
				"/bin",
				"/lib",
				"/opt",
				"/etc",
				"/usr",
				"/tmp",
				"/var",
				"/run",
				"/mnt",
			},
		},
		"name": "PapayaNet OS",
	}

	fmt.Println(data["root"])
	fmt.Println(data)
	fmt.Println("result", koala.KMapGetValue("root.systems.0", data))
	fmt.Println("result", koala.KMapSetValue("root.systems.0", "/usr/bin", data))
	fmt.Println("result", koala.KMapGetValue("root.systems.0", data))
	fmt.Println("result", koala.KMapGetValue("name", data))
	fmt.Println("result", koala.KMapSetValue("name", "Orion X", data))
	fmt.Println("result", koala.KMapGetValue("name", data))
	fmt.Println("result", koala.KMapDelValue("root.systems.1", data))
	fmt.Println(data["root"])
	fmt.Println(data)
}
