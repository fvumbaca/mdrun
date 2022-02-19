package rundoc

import (
	"crypto/sha256"
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

func (d *Rundoc) WriteHTML(w io.Writer) {
	r := customRenderer{
		jsAppURL: "/-/static/js/app.js",
		cssURLs:  []string{"/-/static/css/normalize.css"},
		HTMLRenderer: blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
			CSS:   "/-/static/css/normalize.css",
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
	cssURLs  []string
	*blackfriday.HTMLRenderer
	walkCounter int
}

func (r *customRenderer) RenderHeader(w io.Writer, ast *blackfriday.Node) {
	// r.HTMLRenderer.RenderHeader(w, ast)
	fmt.Fprintln(w, "<!DOCTYPE html>")
	fmt.Fprintln(w, "<html>")
	fmt.Fprintln(w, "<head>")
	fmt.Fprintln(w, "  <title></title>")
	fmt.Fprintln(w, "  <meta name=\"GENERATOR\" content=\"Mdrun powered by Blackfriday 2\">")
	fmt.Fprintln(w, "  <meta charset=\"utf-8\">")
	for _, css := range r.cssURLs {
		fmt.Fprintf(w, "<link rel=\"stylesheet\" href=\"%s\">\n", css)
	}
	fmt.Fprintln(w, "</head>")
	fmt.Fprintln(w, "<body>")
	fmt.Fprintf(w, "<script src='%s'></script>\n", r.jsAppURL)
}

func (r *customRenderer) RenderNode(w io.Writer, node *blackfriday.Node,
	entering bool) blackfriday.WalkStatus {
	r.walkCounter++
	switch node.Type {
	case blackfriday.CodeBlock:
		res := r.HTMLRenderer.RenderNode(w, node, entering)
		block := codeBlock{
			Lang:   string(node.CodeBlockData.Info),
			Script: node.Literal,
		}
		bid := block.GenID()

		fmt.Fprintf(w, "<div id=\"%s\"><button onclick=\"execBlock('%s')\">Run</button></div>", bid, bid)

		// fmt.Fprintf(w, `<div id="block-%d"><button onclick="execBlock('block-%d', '%s')">Run</button></div>`, r.walkCounter, r.walkCounter, block.GenID())
		return res
	default:
		return r.HTMLRenderer.RenderNode(w, node, entering)
	}
}
