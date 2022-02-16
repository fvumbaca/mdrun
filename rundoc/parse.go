package rundoc

import (
	"crypto/md5"
	_ "embed"
	"encoding/base64"
	"fmt"
	"io"

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

func (b codeBlock) GenID() string {
	return base64.StdEncoding.EncodeToString(md5.New().Sum([]byte(b.Lang +
		string(b.Script))))
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

func (d *Rundoc) WriteHTML(w io.Writer) {
	r := customRenderer{
		jsAppURL: "/-/static/js/app.js",
		HTMLRenderer: blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
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
	jsAppURL string
	*blackfriday.HTMLRenderer
}

func (r *customRenderer) RenderHeader(w io.Writer, ast *blackfriday.Node) {
	r.HTMLRenderer.RenderHeader(w, ast)
	fmt.Fprintf(w, "<script src='%s'></script>\n", r.jsAppURL)
}

func (r *customRenderer) RenderNode(w io.Writer, node *blackfriday.Node,
	entering bool) blackfriday.WalkStatus {
	switch node.Type {
	case blackfriday.CodeBlock:
		r := r.HTMLRenderer.RenderNode(w, node, entering)
		block := codeBlock{
			Lang:   string(node.CodeBlockData.Info),
			Script: node.Literal,
		}
		fmt.Fprintf(w, `<button onclick="execBlock('%s')">Run</button>`, block.GenID())
		return r
	default:
		return r.HTMLRenderer.RenderNode(w, node, entering)
	}
}
