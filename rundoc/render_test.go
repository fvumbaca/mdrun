package rundoc

import (
	"bytes"
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "Update golden file for html building tests")

func TestRenderHTML(t *testing.T) {
	var buff bytes.Buffer

	doc, err := Parse(loadFixture(t, "basic.md"))
	noErr(t, err)
	RenderDoc(&buff, doc)

	goldenFilename := filepath.Join("fixtures", t.Name()+".golden.html")
	if *update {
		noErr(t, ioutil.WriteFile(goldenFilename, buff.Bytes(), 0664))
	}
	golden, err := ioutil.ReadFile(goldenFilename)

	noErr(t, err)

	diff(t, string(golden), buff.String())
}
