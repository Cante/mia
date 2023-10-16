package archive

import (
	"os"
	"path/filepath"
)

func GetFileFromDirector(dir string, recursive bool) ([]string, error) {
	var files []string

	info, err := os.Stat(dir)

	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return append(files, dir), nil
	}

	filesInDir, err := os.ReadDir(dir)

	if err != nil {
		return nil, err
	}

	for _, f := range filesInDir {
		if f.IsDir() && recursive {
			ff, err := GetFileFromDirector(filepath.Join(dir, f.Name()), true)

			if err != nil {
				continue
			}

			files = append(files, ff...)

			continue
		}

		files = append(files, filepath.Join(dir, f.Name()))
	}

	return files, nil
}

func GetAccumulatedFileSize(files []string) int64 {
	size := int64(0)

	for _, f := range files {
		info, err := os.Stat(f)

		if err != nil {
			continue
		}

		size += info.Size()
	}

	return size
}
