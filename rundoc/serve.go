package rundoc

import (
	"bytes"
	"context"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type CMDFunc func(context.Context, []byte) ([]byte, error)

type LangExecCMDMap map[string]CMDFunc

var DefaultCMDMap = LangExecCMDMap{
	// "sh": "sh -c",
	"sh": buildStdinCMDFunc("sh"),
	"text": func(ctx context.Context, script []byte) ([]byte, error) {
		return script, nil
	},
}

func buildStdinCMDFunc(bin string) CMDFunc {
	return func(ctx context.Context, script []byte) ([]byte, error) {
		var buff bytes.Buffer
		c := exec.CommandContext(ctx, bin)
		c.Stdin = bytes.NewReader(script)
		c.Env = os.Environ()
		c.Stdout = &buff
		c.Stderr = &buff
		err := c.Run()
		return buff.Bytes(), err
	}
}

type Handler struct {
	RootFS fs.FS
	cmdMap LangExecCMDMap
}

func NewHandler(rootDir fs.FS) *Handler {
	return &Handler{
		RootFS: rootDir,
		cmdMap: DefaultCMDMap,
	}
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
			c, handle := h.cmdMap[b.Lang]
			if !handle {
				// FIXME: Handle the error by responding with a 400
				res, _ := buildStdinCMDFunc(b.Lang)(req.Context(), b.Script)
				w.Write(res)
			} else {
				// FIXME: Handle the error by responding with a 400
				res, _ := c(req.Context(), b.Script)
				w.Write(res)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		doc.WriteHTML(w)
	}
}
