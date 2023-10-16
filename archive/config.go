package archive

import (
	"os"
	"path/filepath"
)

const (
	ModeCompress = 'c'
	ModeExtract  = 'e'
)

type Config struct {
	Target string
	Inputs []string

	Mode       int32
	Recursive  bool
	WorkingDir string
}

func GetConfig() (Config, error) {
	config := Config{}

	workingDir, err := os.MkdirTemp(os.TempDir(), "")
	workingDir = filepath.Join("C:\\xampp\\htdocs\\mia\\_tmp")
	if err != nil {
		return config, nil
	}

	config.WorkingDir = workingDir

	arguments := os.Args[1:]

	for i := 0; i < len(arguments); i++ {
		val := arguments[i]

		switch val {
		case "-c":
			config.Mode = ModeCompress
		case "-e":
			config.Mode = ModeExtract
		case "-r":
			config.Recursive = true
		case "-t":
		case "-o":
			i++
			target, err := filepath.Abs(arguments[i])

			if err != nil {
				return Config{}, err
			}

			config.Target = target
		default:
			input, err := filepath.Abs(val)
			if err == nil {
				if config.Target == "" {
					info, err := os.Stat(input)

					if err != nil {
						return Config{}, err
					}

					if !info.IsDir() {
						input = filepath.Dir(input)
					}

					config.Target = filepath.Join(filepath.Dir(input), filepath.Base(input)+".mia")
				}

				config.Inputs = append(config.Inputs, input)
			}
		}
	}

	return config, nil
}
