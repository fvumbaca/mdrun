package rundoc

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestRenderHTML(t *testing.T) {
	var buff bytes.Buffer
	r := NewHTMLRenderer()

	doc, err := Parse(loadFixture(t, "basic.md"))
	noErr(t, err)
	r.RenderHeader(&buff, doc.docRoot)
	r.RenderFooter(&buff, doc.docRoot)

	goldenFilename := filepath.Join("fixtures", t.Name()+".golden.html")
	if *update {
		noErr(t, ioutil.WriteFile(goldenFilename, buff.Bytes(), 0664))
	}
	golden, err := ioutil.ReadFile(goldenFilename)

	noErr(t, err)

	diff(t, string(golden), buff.String())
}
