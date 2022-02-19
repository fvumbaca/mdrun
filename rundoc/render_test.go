package rundoc

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	blackfriday "github.com/russross/blackfriday/v2"
)

func TestRenderHTML(t *testing.T) {
	var buff bytes.Buffer
	r := NewHTMLRenderer()

	doc, err := Parse(loadFixture(t, "basic.md"))
	noErr(t, err)
	r.RenderHeader(&buff, doc.docRoot)
	doc.docRoot.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		return r.RenderNode(&buff, node, entering)
	})
	r.RenderFooter(&buff, doc.docRoot)

	goldenFilename := filepath.Join("fixtures", t.Name()+".golden.html")
	if *update {
		noErr(t, ioutil.WriteFile(goldenFilename, buff.Bytes(), 0664))
	}
	golden, err := ioutil.ReadFile(goldenFilename)

	noErr(t, err)

	diff(t, string(golden), buff.String())
}