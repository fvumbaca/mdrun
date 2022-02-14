package rundoc

import (
	"fmt"
	"io"

	blackfriday "github.com/russross/blackfriday/v2"
)

type Rundoc struct {
	docRoot *blackfriday.Node
}

func Parse(input []byte) (*Rundoc, error) {
	var doc Rundoc
	optList := []blackfriday.Option{
		blackfriday.WithExtensions(blackfriday.CommonExtensions),
	}
	markdown := blackfriday.New(optList...)
	doc.docRoot = markdown.Parse(input)
	return &doc, nil
}

func (d *Rundoc) WriteHTML(w io.Writer) {
	r := customRenderer{
		blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
			Flags: blackfriday.CompletePage,
		}),
	}
	r.RenderHeader(w, d.docRoot)
	d.docRoot.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		return r.RenderNode(w, node, entering)
	})
	r.RenderFooter(w, d.docRoot)
}

type customRenderer struct {
	*blackfriday.HTMLRenderer
}

func (r *customRenderer) RenderHeader(w io.Writer, ast *blackfriday.Node) {
	r.HTMLRenderer.RenderHeader(w, ast)
}

func (r *customRenderer) RenderNode(w io.Writer, node *blackfriday.Node,
	entering bool) blackfriday.WalkStatus {
	switch node.Type {
	case blackfriday.CodeBlock:
		r := r.HTMLRenderer.RenderNode(w, node, entering)
		fmt.Fprintf(w, `<button onclick="alert('testing123')">Run</button>`)
		return r
	default:
		return r.HTMLRenderer.RenderNode(w, node, entering)
	}
}
