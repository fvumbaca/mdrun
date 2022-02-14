package rundoc

import (
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
	r := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{})
	r.RenderHeader(w, d.docRoot)
	d.docRoot.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		return r.RenderNode(w, node, entering)
	})
	r.RenderFooter(w, d.docRoot)
}
