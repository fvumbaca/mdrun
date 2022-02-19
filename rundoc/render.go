package rundoc

import (
	"embed"
	"html/template"
	"io"

	blackfriday "github.com/russross/blackfriday/v2"
)

//go:embed templates/*
var templatesFS embed.FS
var defaultTemplates *template.Template

func init() {
	var err error
	defaultTemplates, err = template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		panic(err)
	}
}

type HTMLRenderer struct {
	templates    *template.Template
	baseRenderer *blackfriday.HTMLRenderer
	info         renderInfo
}

type renderInfo struct {
	Title    string
	Filepath []string
	CSS      []string
	JS       []string
}

func NewHTMLRenderer() *HTMLRenderer {
	return &HTMLRenderer{
		baseRenderer: blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{}),
		templates:    defaultTemplates,
		info: renderInfo{
			Title:    "mdrun doc",
			Filepath: nil,
			JS:       []string{"/-/static/js/app.js"},
			CSS: []string{
				"/-/static/css/normalize.js",
				"/-/static/css/style.js",
			},
		},
	}
}

func (r *HTMLRenderer) RenderHeader(w io.Writer, ast *blackfriday.Node) {
	r.templates.ExecuteTemplate(w, "header.html", r.info)
}

func (r *HTMLRenderer) RenderFooter(w io.Writer, ast *blackfriday.Node) {
	r.templates.ExecuteTemplate(w, "footer.html", r.info)
}

func (r *HTMLRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	switch node.Type {
	default:
		return r.baseRenderer.RenderNode(w, node, entering)
	}
}
