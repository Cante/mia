package archive

import (
	"compress/zlib"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/barasher/go-exiftool"
	"github.com/schollz/progressbar/v3"
	"image/jpeg"
	"mia/archive/entry"
	"os"
	"path/filepath"
)

var Images = ImageList{}
var Blocks = BlockList{}

type BlockFile struct {
	Name string
	Data []byte
}

func Compress(config Config, inputs []string) error {
	Images = ImageList{
		Images: make([]entry.Image, len(inputs)),
	}

	Blocks = BlockList{
		Blocks: make(map[entry.Hash]entry.SimpleBlock),
	}

	wd, _ := os.Getwd()

	e, err := exiftool.NewExiftool(exiftool.SetExiftoolBinaryPath(filepath.Join(wd, "lib", "exiftool.exe")))
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(inputs)), fmt.Sprintf("Process file %d from %d", 1, len(inputs)))

	for i, file := range inputs {
		handle, _ := os.Open(file)

		img, _ := jpeg.Decode(handle)
		header, _ := jpeg.DecodeConfig(handle)

		fileInfos := e.ExtractMetadata(file)
		j, err := json.Marshal(fileInfos[0].Fields)

		if err != nil {
			return err
		}

		image := entry.Image{
			ImageFile:   img,
			ImageHeader: header,
			Exif:        j,
		}

		blocks := image.ToBlocks(true)

		for _, block := range blocks {
			image.BlockHashes = append(image.BlockHashes, block.Hash)

			if Blocks.HasBlock(block.Hash) {
				continue
			}

			_, err := SaveBlock(config, block)

			if err != nil {
				return err
			}

			Blocks.AddBlock(block)
		}

		Images.Add(image)

		err = handle.Close()
		if err != nil {
			return err
		}

		bar.Describe(fmt.Sprintf("Process file %d from %d", i+1, len(inputs)))

		if i+1 == len(inputs) {
			bar.Clear()
			break
		}

		bar.Add(1)
	}

	fmt.Println(fmt.Sprintf("Processing %d files:\tFinish", len(inputs)))

	return saveMiaFile(config)
}

//goland:noinspection GoUnhandledErrorResult
func saveMiaFile(config Config) error {
	file, err := os.Create(config.Target)
	if err != nil {
		return err
	}

	defer file.Close()

	zlibWriter := zlib.NewWriter(file)
	defer zlibWriter.Close()

	encoder := gob.NewEncoder(zlibWriter)
	if err := encoder.Encode(Images); err != nil {
		return err
	}

	if err := encoder.Encode(Blocks); err != nil {
		return err
	}

	blockFiles, _ := os.ReadDir(config.WorkingDir)

	bar := progressbar.Default(int64(len(blockFiles)), "Compressing")

	for i, blockFile := range blockFiles {
		path := filepath.Join(config.WorkingDir, blockFile.Name())
		data, err := os.ReadFile(path)

		if err != nil {
			return err
		}

		err = encoder.Encode(BlockFile{
			Name: blockFile.Name(),
			Data: data,
		})

		if err != nil {
			return err
		}

		if i+1 == len(blockFiles) {
			bar.Clear()
			break
		}

		bar.Add(1)
	}

	fmt.Println(fmt.Sprintf("Compression:\t\tFinish"))

	return nil
}
