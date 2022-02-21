package documents

import (
	"io"
	"strings"

	bf "github.com/russross/blackfriday/v2"
)

func RenderHTML(w io.Writer, doc *Document) {
	doc.root.Walk(func(node *bf.Node, entering bool) bf.WalkStatus {
		return renderNode(w, node, entering)
	})
}

var tags = map[bf.NodeType][]string{
	bf.Document:       {"", ""},
	bf.BlockQuote:     {"", ""},
	bf.List:           {"", ""},
	bf.Item:           {"", ""},
	bf.Paragraph:      {"", ""},
	bf.Heading:        {"", ""},
	bf.HorizontalRule: {"", ""},
	bf.Emph:           {"", ""},
	bf.Strong:         {"", ""},
	bf.Del:            {"", ""},
	bf.Link:           {"", ""},
	bf.Image:          {"", ""},
	bf.Text:           {"", ""},
	bf.HTMLBlock:      {"", ""},
	bf.CodeBlock:      {"", ""},
	bf.Softbreak:      {"", ""},
	bf.Hardbreak:      {"", ""},
	bf.Code:           {"", ""},
	bf.HTMLSpan:       {"", ""},
	bf.Table:          {"", ""},
	bf.TableCell:      {"", ""},
	bf.TableHead:      {"", ""},
	bf.TableBody:      {"", ""},
	bf.TableRow:       {"", ""},
}

func renderNode(w io.Writer, node *bf.Node, entering bool) bf.WalkStatus {
	var ident = 0
	io.WriteString(w, strings.Repeat("  ", ident))
	switch node.Type {
	default:
		if entering {
			io.WriteString(w, tags[node.Type][0])
		} else {
			io.WriteString(w, tags[node.Type][1])
		}
	}
	if entering {
		ident++
	} else {
		ident--
	}
	io.WriteString(w, "\n")
	return bf.GoToNext
}
