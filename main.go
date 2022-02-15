package main

/*
This is going to be sick!

Features:

- Support metadata at the top of the file
- Support required cli-tools to throw errors early if it can even run on the
system
- Support inline requests
- Support out-of-the-box containerization
- Interactive mode will print out runbook and prompt before each exec
- Script mode that will execute mardown as a script
- Infer language dependencies from code blocks

Example usage:

	$ mdrun my-runbook.md --set SOME_VAR=hello

	# -E pulls from the environment
	$ mdrun -E my-runbook.md
*/

import (
	// Keeping here for now to include tests while developing

	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fvumbaca/mdrun/browser"
	_ "github.com/fvumbaca/mdrun/rundoc"
	"github.com/gorilla/websocket"
	"github.com/russross/blackfriday/v2"
)

func main() {
	fmt.Println("Startung up server on :3000")
	http.Handle("/", &browser.Handler{RootFS: os.DirFS("./")})
	http.ListenAndServe(":3000", nil)
}

func xmain() {
	// http.HandleFunc("/", listFilesHandler)
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/-/ws", handleWs)
	http.Handle("/", &Handler{fs: os.DirFS(cwd)})
	http.ListenAndServe(":3000", nil)
}

var pageListT *template.Template

func init() {
	var err error
	pageListT, err = template.New("").Parse(`
<h1>Files</h1>
<ul>
{{ range . }}
<li><a href="{{.}}">{{ . }}</a></li>
{{ end }}
</ul>
	`)
	if err != nil {
		panic(err)
	}
}

type Handler struct {
	fs fs.FS
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// p := strings.TrimPrefix(req.URL.Path, "/")
	// p := filepath.Clean("." + req.URL.Path)
	p := strings.Trim(req.URL.Path, "/")
	if p == "" {
		p = "."
	}

	f, err := h.fs.Open(p)
	// defer f.Close()
	if err != nil {
		returnError(w, req, err)
		return
	}
	fi, err := f.Stat()
	if err != nil {
		returnError(w, req, err)
		return
	}

	if fi.IsDir() {
		var filenames []string
		fs.WalkDir(h.fs, p, func(pp string, entry fs.DirEntry, err error) error {
			if p == pp {
				return nil
			}

			if entry.IsDir() {
				filenames = append(filenames, entry.Name()+"/")
				return fs.SkipDir
			}
			if filepath.Ext(entry.Name()) == ".md" {
				filenames = append(filenames, entry.Name())
			}
			return nil
		})
		pageListT.Execute(w, filenames)
	} else if strings.HasSuffix(fi.Name(), "md") {
		contents, err := ioutil.ReadAll(f)
		if err != nil {
			returnError(w, req, err)
			return
		}
		var doc Rundoc
		rendered := blackfriday.Run(contents, blackfriday.WithRenderer(NewCodeBlockRenderer(&doc)))
		w.Write(rendered)
		// io.Copy(w, bytes.NewReader(rendered))

	} else {
		returnError(w, req, errors.New("Not found....."))
	}
}

func returnError(w http.ResponseWriter, req *http.Request, err error) {
	fmt.Println("Error when rendering: ", err)
	fmt.Fprintln(w, "Error!")
}

type CodeBlockRenderer struct {
	rundoc *Rundoc
	*blackfriday.HTMLRenderer
}

func NewCodeBlockRenderer(doc *Rundoc) *CodeBlockRenderer {
	return &CodeBlockRenderer{
		rundoc: doc,
		HTMLRenderer: blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
			// TODO
			Title: "Rundoc",
			CSS:   "https://cdn.tailwindcss.com",
			Flags: blackfriday.CompletePage,
		}),
	}
}

func (r *CodeBlockRenderer) RenderHeader(w io.Writer, ast *blackfriday.Node) {
	r.HTMLRenderer.RenderHeader(w, ast)
	r.HTMLRenderer.RenderHeader(os.Stdout, ast)
	fmt.Fprintln(w, `<script>
var s = new WebSocket(((window.location.protocol === "https:") ? "wss://" : "ws://") + window.location.host + "/-/ws");

s.addEventListener('message', function (event) {
    console.log('Message from server ', event);
});

function execBase64(based) {
	// alert("Going to execute " + based);
    s.send(JSON.stringify({
		lang: "sh",
		script: "Hello World!",
	}));

}
</script>`)
}

func (r *CodeBlockRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	switch node.Type {
	case blackfriday.CodeBlock:
		// fmt.Println(string(node.CodeBlockData.Info))
		// fmt.Printf("X %#v\n", node)
		// fmt.Fprintf(w, "<pre><code class=\"language-sh\">echo &quot;hello, world!RENDERED&quot;</code></pre>")
		r.HTMLRenderer.RenderNode(w, node, entering)
		// fmt.Fprintf(w, "<div><button>Run</button>%s</div>\n", string(node.Literal))
		renderPlayButton(w, node.Literal)
		return blackfriday.GoToNext
	default:
		return r.HTMLRenderer.RenderNode(w, node, entering)
	}
}

func renderPlayButton(w io.Writer, script []byte) {
	based := base64.StdEncoding.EncodeToString(script)
	fmt.Fprintf(w, `<button onclick="execBase64('%s')">Run</button>`, based)
}

type Rundoc struct {
	scripts map[string][]byte
}

var upgrader = websocket.Upgrader{}

func handleWs(w http.ResponseWriter, req *http.Request) {
	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		returnError(w, req, err)
		return
	}

	defer ws.Close()

	var msg ExecReq
	for err == nil {
		err = ws.ReadJSON(&msg)
		if err != nil {
			fmt.Println("it happened", err)
			break
		}

		resp := ExecResp{
			Result: "this came from the server",
		}
		err = ws.WriteJSON(&resp)
		if err != nil {
			fmt.Println("it happened", err)
			break
		}
	}
}

type ExecReq struct {
	Lang   string `json:"lang"`
	Script string `json:"script"`
}

type ExecResp struct {
	Result string `json:"result"`
}
