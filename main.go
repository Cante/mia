package main

import (
	"fmt"
	"mia/archive"
)

func main() {
	config, err := archive.GetConfig()

	if err != nil {
		panic(err)
	}

	//defer func() {
	//	err := os.RemoveAll(config.WorkingDir)
	//	if err != nil {
	//		panic(err)
	//	}
	//}()

	var inputs []string

	for _, input := range config.Inputs {
		files, err := archive.GetFileFromDirector(input, true)

		if err != nil {
			continue
		}

		inputs = append(inputs, files...)
	}

	archive.PrintWelcome(config, inputs)
	err = archive.Compress(config, inputs)

	fmt.Println(err)

	if err != nil {
		fmt.Errorf("%v", err)
	}
}
