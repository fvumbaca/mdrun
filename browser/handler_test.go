package browser

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/fvumbaca/mdrun/rundoc"
)

func TestHandler(t *testing.T) {
	dir, cleanup := tmpDir(t)
	defer cleanup()

	newFile(t, dir, "example.md", "")
	newDir(t, dir, "subdir")
	newFile(t, dir, "subdir/subfile.md", "")

	// FIXME: There should be a better way to test this then to import rundoc's
	// handler. Did this bc it works for now....
	fh := rundoc.Handler{RootFS: os.DirFS(dir)}
	h := &Handler{
		RootFS:      os.DirFS(dir),
		FileHandler: &fh,
	}
	hitHTMLEndpointStatusCode(t, h, "/does_not_exist.md", http.StatusNotFound)
	hitHTMLEndpointStatusCode(t, h, "/example.md", http.StatusOK)
	hitHTMLEndpointStatusCode(t, h, "/subdir", http.StatusOK)
	hitHTMLEndpointStatusCode(t, h, "/subdir_not_exist", http.StatusNotFound)
	hitHTMLEndpointStatusCode(t, h, "/", http.StatusOK)
	hitHTMLEndpointStatusCode(t, h, "/subdir/subfile.md", http.StatusOK)
}

func hitHTMLEndpointStatusCode(t *testing.T, h http.Handler, endpoint string, expected int) {
	r := httptest.NewRequest("GET", endpoint, nil)
	rq := httptest.NewRecorder()
	h.ServeHTTP(rq, r)

	res := rq.Result()
	if expected != res.StatusCode {
		t.Error("GET", endpoint, "expected to return", expected, "but got",
			res.StatusCode, "instead")
	}
}
