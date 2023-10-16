package entry

import (
	"github.com/cespare/xxhash/v2"
	"image"
	"image/color"
)

const LineSize = 256
const BlockSize = LineSize * LineSize

type Hash uint64
type Pixel struct {
	RGBA  color.RGBA
	valid bool
}

func (p *Pixel) Valid() bool {
	return p.valid
}

func (p *Pixel) ToByte() []byte {
	return []byte{
		p.RGBA.R,
		p.RGBA.G,
		p.RGBA.B,
		p.RGBA.A,
	}
}

type Block struct {
	Hash  Hash
	Pixel [BlockSize]Pixel
	Size  struct {
		X int
		Y int
	}
}

type SimpleBlock struct {
	Hash Hash
	Size struct {
		X int
		Y int
	}
}

func (b *Block) AsSimpleBlock() SimpleBlock {
	return SimpleBlock{
		Hash: b.Hash,
		Size: b.Size,
	}
}

func (b *Block) hash() {
	var bytes []byte

	for _, p := range b.Pixel {
		bytes = append(bytes, p.ToByte()...)
	}

	bytes = append(bytes, byte(b.Size.X), byte(b.Size.Y))

	b.Hash = Hash(xxhash.Sum64(bytes))
}

type Image struct {
	ImageFile   image.Image
	ImageHeader image.Config
	Exif        []byte
	BlockHashes []Hash

	meta struct {
		width  int
		height int
	}
}

func (i *Image) ToBlocks(minimize bool) []Block {
	bounds := i.ImageFile.Bounds()
	i.meta.width, i.meta.height = bounds.Max.X, bounds.Max.Y

	var blocks []Block

	for y := 0; y < i.meta.height; y += LineSize {
		for x := 0; x < i.meta.width; x += LineSize {
			block := i.extractBlock(x, y, LineSize)
			blocks = append(blocks, block)
		}
	}

	if minimize {
		i.ImageFile = nil
	}

	return blocks
}

func (i *Image) extractBlock(x, y, size int) Block {
	bounds := i.ImageFile.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	block := Block{}

	for yOffset := 0; yOffset < size; yOffset++ {
		for xOffset := 0; xOffset < size; xOffset++ {
			if x+xOffset < width && y+yOffset < height {
				p := color.RGBAModel.Convert(i.ImageFile.At(x+xOffset, y+yOffset)).(color.RGBA)
				block.Pixel[xOffset*LineSize+yOffset] = Pixel{p, true}
				block.Size.X = xOffset
				block.Size.Y = yOffset
			} else {
				block.Pixel[xOffset*LineSize+yOffset] = Pixel{color.RGBA{}, false}
			}
		}
	}

	block.hash()
	i.BlockHashes = append(i.BlockHashes, block.Hash)

	return block
}
