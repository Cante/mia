package archive

import "mia/archive/entry"

type BlockList struct {
	Blocks map[entry.Hash]entry.SimpleBlock
}

func (l *BlockList) Add(block entry.SimpleBlock) {
	if _, ok := l.Blocks[block.Hash]; ok {
		return
	}

	l.Blocks[block.Hash] = block
}

func (l *BlockList) AddBlock(block entry.Block) {
	l.Add(block.AsSimpleBlock())
}

func (l *BlockList) HasBlock(hash entry.Hash) bool {
	if _, ok := l.Blocks[hash]; ok {
		return true
	}

	return false
}

type ImageList struct {
	Images []entry.Image
}

func (l *ImageList) Add(image entry.Image) {
	l.Images = append(l.Images, image)
}
