package rundoc

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
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
