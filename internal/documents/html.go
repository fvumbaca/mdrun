package documents

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	bf "github.com/russross/blackfriday/v2"
)

func RenderHTML(w io.Writer, doc *Document) {
	doc.root.Walk(func(node *bf.Node, entering bool) bf.WalkStatus {
		return renderNode(w, node, entering)
	})
}

var tags = map[bf.NodeType][]string{
	// bf.List:           {"", ""},
	// bf.Item:           {"", ""},
	// bf.HorizontalRule: {"", ""},
	// bf.Emph:           {"", ""},
	// bf.Strong:         {"", ""},
	// bf.Del:            {"", ""},
	// bf.Link:           {"", ""},
	// bf.Image:          {"", ""},
	// bf.HTMLBlock:      {"", ""},
	// bf.Softbreak:      {"", ""},
	// bf.Hardbreak:      {"", ""},
	// bf.Code:           {"", ""},
	// bf.HTMLSpan:       {"", ""},
	// bf.Table:          {"", ""},
	// bf.TableCell:      {"", ""},
	// bf.TableHead:      {"", ""},
	// bf.TableBody:      {"", ""},
	// bf.TableRow:       {"", ""},
}

func renderNode(w io.Writer, node *bf.Node, entering bool) bf.WalkStatus {
	var ident = 0
	io.WriteString(w, strings.Repeat("  ", ident))
	switch node.Type {
	case bf.Document:
		return bf.GoToNext
	case bf.Text:
		w.Write(node.Literal)
	case bf.Heading:
		io.WriteString(w, "<")
		if !entering {
			io.WriteString(w, "/")
		}
		io.WriteString(w, "h"+strconv.Itoa(node.HeadingData.Level)+">")
	case bf.Paragraph:
		if entering {
			io.WriteString(w, "<p>")
		} else {
			io.WriteString(w, "</p>")
		}
	case bf.CodeBlock:
		fmt.Fprintf(w, "<pre><code class=\"language-%s\">%s</code></pre>", string(node.CodeBlockData.Info), node.Literal)
	case bf.BlockQuote:
		if entering {
			io.WriteString(w, "<div class=\"well\">")
		} else {
			io.WriteString(w, "</div>")
		}

	default:
		panic(node.Type.String())
	}
	if entering {
		ident++
	} else {
		ident--
	}
	io.WriteString(w, "\n")
	return bf.GoToNext
}
