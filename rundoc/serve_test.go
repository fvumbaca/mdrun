package rundoc

import (
	"bytes"
	"fmt"
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
	h := NewHandler(os.DirFS(dir))
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

	h := NewHandler(os.DirFS(dir))
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

	h := NewHandler(os.DirFS(dir))
	h.ServeHTTP(rec, req)
	r := rec.Result()

	if r.StatusCode != http.StatusNotFound {
		t.Error("expected status 404 but got", r.StatusCode)
	}
}

func TestHandlerExec_TextBlock(t *testing.T) {
	dir, cleanup := tmpDir(t)
	defer cleanup()
	newFile(t, dir, "example.md", "example md file\n```text\nExample\n```")
	h := NewHandler(os.DirFS(dir))
	testBlockExecEndpoint(t, h, "example.md", "text", "Example\n", "Example\n")
}

func TestHandlerExec_ShellEchoBlock(t *testing.T) {
	dir, cleanup := tmpDir(t)
	defer cleanup()
	newDocBuilder().AddCodeBlock(codeBlock{
		Lang:   "sh",
		Script: []byte("echo Hello world"),
	}).WriteFile(t, dir, "t.md")

	h := NewHandler(os.DirFS(dir))
	testBlockExecEndpoint(t, h, "t.md", "sh", "echo \"hello, world!\"\n", "Hello world\n")
}

func TestHandlerExec_ShellNotInCMDFuncMap(t *testing.T) {
	dir, cleanup := tmpDir(t)
	defer cleanup()
	newDocBuilder().AddCodeBlock(codeBlock{
		Lang:   "sh",
		Script: []byte("echo Hello world"),
	}).WriteFile(t, dir, "t.md")

	h := NewHandler(os.DirFS(dir))
	h.cmdMap = make(LangExecCMDMap)
	testBlockExecEndpoint(t, h, "t.md", "sh", "echo \"hello, world!\"\n", "Hello world\n")
}

// helpers ---------

type docBuilder struct {
	bytes.Buffer
	blocks []codeBlock
}

func newDocBuilder() *docBuilder {
	return &docBuilder{}
}

func (b *docBuilder) AddCodeBlock(block codeBlock) *docBuilder {
	b.blocks = append(b.blocks, block)
	fmt.Fprintf(b, "```%s\n%s\n```\n", block.Lang, string(block.Script))
	return b
}

func (b *docBuilder) WriteFile(t *testing.T, dir, filename string) {
	newFile(t, dir, filename, b.String())
}

func testBlockExecEndpoint(t *testing.T, h http.Handler, filename, lang, script, expected string) {
	bid := codeBlock{
		Lang:   lang,
		Script: []byte(script),
	}.GenID()
	req := httptest.NewRequest("POST", "/"+filename+"?bid="+bid, nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)
	r := rec.Result()

	if r.StatusCode != http.StatusOK {
		t.Error("expected status 200 but got", r.StatusCode)
	}

	content, err := ioutil.ReadAll(r.Body)
	noErr(t, err)
	diff(t, expected, string(content))
}

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
