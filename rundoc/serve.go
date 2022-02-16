package rundoc

import (
	"io/fs"
	"net/http"
	"strings"
)

type Handler struct {
	RootFS fs.FS
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	filename := strings.TrimPrefix(req.URL.Path, "/")

	contents, err := fs.ReadFile(h.RootFS, filename)
	if perr, ok := err.(*fs.PathError); ok && perr.Op == "open" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	doc, err := Parse(contents)
	if err != nil {
		// TODO: Render a nice error page when there is an issue building from
		// the md file
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Method == http.MethodPost {
		b, ok := doc.blocks[req.URL.Query().Get("bid")]
		if ok {
			w.Write(b.Script)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		doc.WriteHTML(w)
	}
}
