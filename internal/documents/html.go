package documents

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	bf "github.com/russross/blackfriday/v2"
)

func RenderHTML(w io.Writer, doc *Document) {
	r := renderer{}
	doc.root.Walk(func(node *bf.Node, entering bool) bf.WalkStatus {
		return r.renderNode(w, node, entering)
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

type renderer struct {
	isListItemText bool
}

func (r *renderer) renderNode(w io.Writer, node *bf.Node, entering bool) bf.WalkStatus {
	var ident = 0
	io.WriteString(w, strings.Repeat("  ", ident))
	switch node.Type {
	case bf.HTMLBlock:
		w.Write(node.Literal)
	case bf.Document:
		return bf.GoToNext

	case bf.Text:
		txt := string(node.Literal)
		if r.isListItemText {
			if strings.HasPrefix(txt, "[ ] ") {
				io.WriteString(w, "<input type=\"checkbox\" disabled=\"\">")
				txt = strings.TrimPrefix(txt, "[ ] ")
			}
			if strings.HasPrefix(txt, "[x] ") {
				io.WriteString(w, "<input type=\"checkbox\" checked=\"\" disabled=\"\">")
				txt = strings.TrimPrefix(txt, "[x] ")
			}
		}
		io.WriteString(w, txt)
		r.isListItemText = false
	case bf.Heading:
		// TODO: add anchors
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
	case bf.Link:
		// TODO: support anchors
		if entering {
			fmt.Fprintf(w, "<a href=\"%s\">", string(node.LinkData.Destination))
		} else {
			io.WriteString(w, "</a>")
		}
	case bf.Strong:
		if entering {
			io.WriteString(w, "<strong>")
		} else {
			io.WriteString(w, "</strong>")
		}
	case bf.Emph:
		if entering {
			io.WriteString(w, "<em>")
		} else {
			io.WriteString(w, "</em>")
		}
	case bf.Code:
		if entering {
			io.WriteString(w, "<code>")
		} else {
			io.WriteString(w, "</code>")
		}
	case bf.Image:
		if entering {
			fmt.Fprintf(w, "<img src=\"%s\">", string(node.LinkData.Destination))
		} else {
			io.WriteString(w, "</img>")
		}
	case bf.List:
		if node.ListData.ListFlags&bf.ListTypeOrdered != 0 {
			if entering {
				io.WriteString(w, "<ol>")
			} else {
				io.WriteString(w, "</ol>")
			}
		} else {
			if entering {
				io.WriteString(w, "<ul>")
			} else {
				io.WriteString(w, "</ul>")
			}
		}
	case bf.Item:
		if entering {
			io.WriteString(w, "<li>")
			r.isListItemText = true
		} else {
			io.WriteString(w, "</li>")
			r.isListItemText = false
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
