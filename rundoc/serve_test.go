package rundoc

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestHandler_FileNotFound(t *testing.T) {
	dir, cleanup := tmpDir(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/example.md", nil)
	rec := httptest.NewRecorder()
	h := Handler{os.DirFS(dir)}
	h.ServeHTTP(rec, req)
	r := rec.Result()

	if r.StatusCode != http.StatusNotFound {
		t.Error("expected status 404 but got", r.StatusCode)
	}
}

func TestHandler_ExistingFile(t *testing.T) {
	dir, cleanup := tmpDir(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/example.md", nil)
	rec := httptest.NewRecorder()

	newFile(t, dir, "example.md", "example md file")

	h := Handler{os.DirFS(dir)}
	h.ServeHTTP(rec, req)
	r := rec.Result()

	if r.StatusCode != http.StatusOK {
		t.Error("expected status 200 but got", r.StatusCode)
	}
}

func TestHandler_ExecBlockDoesNotExist(t *testing.T) {
	dir, cleanup := tmpDir(t)
	defer cleanup()

	req := httptest.NewRequest("POST", "/example.md?bid=123", nil)
	rec := httptest.NewRecorder()

	newFile(t, dir, "example.md", "example md file")

	h := Handler{os.DirFS(dir)}
	h.ServeHTTP(rec, req)
	r := rec.Result()

	if r.StatusCode != http.StatusNotFound {
		t.Error("expected status 404 but got", r.StatusCode)
	}
}

func TestHandler_ExecBlock(t *testing.T) {
	dir, cleanup := tmpDir(t)
	defer cleanup()

	newFile(t, dir, "example.md", "example md file\n```sh\nExample\n```")

	bid := codeBlock{
		Lang:   "sh",
		Script: []byte("Example\n"),
	}.GenID()

	req := httptest.NewRequest("POST", "/example.md?bid="+bid, nil)
	rec := httptest.NewRecorder()

	h := Handler{os.DirFS(dir)}
	h.ServeHTTP(rec, req)
	r := rec.Result()

	if r.StatusCode != http.StatusOK {
		t.Error("expected status 200 but got", r.StatusCode)
	}

	content, err := ioutil.ReadAll(r.Body)
	noErr(t, err)

	diff(t, string(content), "Example\n")
}

// helpers ---------

func newFile(t *testing.T, dir, name string, body string) {
	var err error
	if body != "" {
		err = ioutil.WriteFile(filepath.Join(dir, name), []byte(body), 0655)
	} else {
		var f *os.File
		f, err = os.Create(filepath.Join(dir, name))
		defer f.Close()
	}
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func newDir(t *testing.T, dir, name string) {
	err := os.MkdirAll(filepath.Join(dir, name), 0775)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func diff(t *testing.T, exp, act interface{}) {
	if diff := cmp.Diff(exp, act); diff != "" {
		t.Error(diff)
	}
}

func noErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func tmpDir(t *testing.T) (string, func()) {
	f, err := os.MkdirTemp("", t.Name())
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	return f, func() {
		os.RemoveAll(f)
	}
}
