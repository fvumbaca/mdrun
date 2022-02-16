package browser

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
)

//go:embed *.tmpl.html
var templateFS embed.FS

var templates *template.Template

func init() {
	var err error
	templates, err = template.ParseFS(templateFS, "*")
	if err != nil {
		panic(err)
	}
}

type Handler struct {
	RootFS      fs.FS
	FileHandler http.Handler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	filename := strings.Trim(req.URL.Path, "/")

	// Handle requests to /
	if filename == "" {
		filename = "."
	}

	fi, err := fs.Stat(h.RootFS, filename)
	if perr, ok := err.(*fs.PathError); ok && perr.Op == "open" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !fi.IsDir() {
		h.FileHandler.ServeHTTP(w, req)
		return
	}
	h.renderDirectory(w, req, filename)
}

func (h *Handler) renderDirectory(w http.ResponseWriter, req *http.Request, dirname string) {
	dir, err := fs.Sub(h.RootFS, dirname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	filenames, err := listMarkdownAndDirs(dir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	templates.ExecuteTemplate(w, "directory.tmpl.html", infoDirectory{
		Name:  filepath.Base(dirname),
		Path:  filepath.Clean(dirname),
		Items: filenames,
	})
}

type infoDirectory struct {
	Name  string
	Path  string
	Items []string
}
