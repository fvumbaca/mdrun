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

func TestRenderGithubDocumented(t *testing.T) {
	var buff bytes.Buffer
	doc := Parse(FixtureBytes(t, "documented-github.md"))
	RenderHTML(&buff, doc)
	GoldenFileDiff(t, "documented_github.golden.html", buff.Bytes())
}
