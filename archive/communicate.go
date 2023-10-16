package archive

import (
	"fmt"
	"mia/util"
	"strconv"
)

func PrintWelcome(config Config, inputs []string) {
	fmt.Println("-----------------------------------------------")
	fmt.Println("---                   Mia                   ---")
	fmt.Println("-----------------------------------------------")

	if config.Mode == ModeCompress {
		fmt.Println("Task:\t\tCompress")
		fmt.Println("Input:\t\t" + strconv.Itoa(len(inputs)) + " files")
		fmt.Println("Size:\t\t" + util.FormatBytes(float64(GetAccumulatedFileSize(inputs))))
		fmt.Println("Output:\t\t" + config.Target)
	} else if config.Mode == ModeExtract {
		fmt.Println("Task:\tExtract")
	}
}
