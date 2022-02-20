package rundoc

import (
	"crypto/sha256"
	_ "embed"
	"encoding/base64"

	blackfriday "github.com/russross/blackfriday/v2"
)

type Rundoc struct {
	docRoot *blackfriday.Node
	blocks  map[string]codeBlock
}

type codeBlock struct {
	Lang   string
	Script []byte
}

func Parse(input []byte) (*Rundoc, error) {
	var doc Rundoc
	optList := []blackfriday.Option{
		blackfriday.WithExtensions(blackfriday.CommonExtensions),
	}
	markdown := blackfriday.New(optList...)
	doc.docRoot = markdown.Parse(input)

	doc.blocks = make(map[string]codeBlock)
	doc.docRoot.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if node.Type == blackfriday.CodeBlock {
			block := codeBlock{
				Lang:   string(node.CodeBlockData.Info),
				Script: node.Literal,
			}
			doc.blocks[block.GenID()] = block
		}
		return blackfriday.GoToNext
	})

	return &doc, nil
}

func (b codeBlock) GenID() string {
	// TODO: Review this to minimize collisions
	id := base64.RawStdEncoding.EncodeToString(sha256.New().Sum([]byte(b.Lang +
		string(b.Script))))
	if len(id) < 8 {
		return id
	}
	return id[:8]
}
