package documents

import (
	"bytes"
	"testing"

	. "github.com/fvumbaca/mdrun/internal/testutil"
)

func TestBasicRender(t *testing.T) {
	var buff bytes.Buffer
	doc := Parse(FixtureBytes(t, "basic.md"))
	RenderHTML(&buff, doc)
	GoldenFileDiff(t, "basic_md_render.golden.html", buff.Bytes())
}
