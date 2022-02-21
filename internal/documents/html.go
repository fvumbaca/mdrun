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

type renderer struct {
	isListItemText bool
}

// TODO: clean up this code. Its really ugly rn....
func (r *renderer) renderNode(w io.Writer, node *bf.Node, entering bool) bf.WalkStatus {
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
		block := CodeBlock{
			Lang:   string(node.CodeBlockData.Info),
			Script: string(node.Literal),
		}
		bid := block.GenID()
		io.WriteString(w, "<div id=\""+bid+"\"><button onclick=\"execBlock('"+bid+"')\">Run</button></div>\n")
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

	case bf.Table:
		if entering {
			io.WriteString(w, "<table>")
		} else {
			io.WriteString(w, "</table>")
		}
	case bf.TableHead:
		if entering {
			io.WriteString(w, "<thead>")
		} else {
			io.WriteString(w, "</thead>")
		}
	case bf.TableRow:
		if entering {
			io.WriteString(w, "<tr>")
		} else {
			io.WriteString(w, "</tr>")
		}
	case bf.TableCell:
		if entering {
			io.WriteString(w, "<td>")
		} else {
			io.WriteString(w, "</td>")
		}
	case bf.TableBody:
		if entering {
			io.WriteString(w, "<tbody>")
		} else {
			io.WriteString(w, "</tbody>")
		}
	case bf.HorizontalRule:
		io.WriteString(w, "</hr>")
	case bf.Hardbreak:
		io.WriteString(w, "</br>")
	case bf.Softbreak:
		io.WriteString(w, "\n")

	default:
		panic(node.Type.String())
	}
	io.WriteString(w, "\n")
	return bf.GoToNext
}
