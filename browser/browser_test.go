package browser

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDirList_NotExist(t *testing.T) {
	filenames, err := listMarkdownAndDirs("/doesnotexist")
	if err == nil {
		t.Error("expected an error but got none")
	}
	diff(t, []string(nil), filenames)
}

func TestDirList(t *testing.T) {
	dir, cleanup := tmpDir(t)
	defer cleanup()

	newFile(t, dir, "example.md", "")
	newDir(t, dir, "subdir")
	newFile(t, dir, "ignoreme.txt", "")

	filenames, err := listMarkdownAndDirs(dir)
	noErr(t, err)

	expected := []string{
		"example.md",
		"subdir",
	}

	if diff := cmp.Diff(expected, filenames); diff != "" {
		t.Error(diff)
	}
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
	err := os.MkdirAll(filepath.Join(dir, name), 0655)
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
