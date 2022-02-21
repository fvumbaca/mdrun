package documents

import (
	"io"

	"github.com/russross/blackfriday/v2"
)

type Document struct {
	blocks map[string]CodeBlock
	root   *blackfriday.Node
}

func Parse(content []byte) Document {
	doc := Document{
		blocks: make(map[string]CodeBlock),
	}
	md := blackfriday.New()
	doc.root = md.Parse(content)
	return doc
}

func testPrintNode(w io.Writer, doc Document) {
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
