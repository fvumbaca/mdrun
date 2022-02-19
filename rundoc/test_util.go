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

func loadFixture(t *testing.T, name string) []byte {
	content, err := ioutil.ReadFile("fixtures/" + name)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	return content
}

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
