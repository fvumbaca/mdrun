package documents

import (
	"io"

	"github.com/russross/blackfriday/v2"
)

type Document struct {
	blocks map[string]CodeBlock
	root   *blackfriday.Node
}

func Parse(content []byte) *Document {
	doc := &Document{
		blocks: make(map[string]CodeBlock),
	}
	md := blackfriday.New(blackfriday.WithExtensions(blackfriday.CommonExtensions))
	doc.root = md.Parse(content)

	doc.root.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if node.Type == blackfriday.CodeBlock {
			block := CodeBlock{
				Lang:   string(node.CodeBlockData.Info),
				Script: string(node.Literal),
			}
			doc.blocks[block.GenID()] = block
		}
		return blackfriday.GoToNext
	})

	return doc
}

func (d *Document) GetCodeBlock(blockID string) (CodeBlock, bool) {
	block, ok := d.blocks[blockID]
	return block, ok
}

func testPrintDocument(w io.Writer, doc *Document) {
	doc.root.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if entering {
			io.WriteString(w, "entering ")
		} else {
			io.WriteString(w, "exiting ")
		}
		io.WriteString(w, node.String()+"\n")
		return blackfriday.GoToNext
	})
}
