package archive

import (
	"compress/zlib"
	imageFile "image"
	"image/jpeg"
	"mia/archive/entry"
	"os"
	"path/filepath"
	"strconv"
)

//goland:noinspection GoUnhandledErrorResult
func SaveBlock(config Config, block entry.Block) (string, error) {
	filename := filepath.Join(config.WorkingDir, strconv.FormatUint(uint64(block.Hash), 10)+".bin")

	if info, _ := os.Stat(filename); info != nil {
		return filename, nil
	}

	img := imageFile.NewRGBA(imageFile.Rect(0, 0, block.Size.X, block.Size.Y))
	for x := 0; x < entry.LineSize; x++ {
		for y := 0; y < entry.LineSize; y++ {
			pixel := block.Pixel[x*entry.LineSize+y]

			if !pixel.Valid() {
				continue
			}

			img.Set(x, y, pixel.RGBA)
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}

	defer file.Close()

	zlibWriter, _ := zlib.NewWriterLevel(file, zlib.BestCompression)
	defer zlibWriter.Close()

	options := jpeg.Options{
		Quality: 100,
	}

	if err := jpeg.Encode(zlibWriter, img, &options); err != nil {
		return "", err
	}

	return filename, nil
}
