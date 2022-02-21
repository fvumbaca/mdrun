package documents

import (
	"bytes"
	"testing"

	. "github.com/fvumbaca/mdrun/internal/testutil"
)

func TestParse_Basic(t *testing.T) {
	doc := Parse(FixtureBytes(t, "basic.md"))
	var buff bytes.Buffer
	testPrintDocument(&buff, doc)
	GoldenFileDiff(t, "basic_parse_result.golden", buff.Bytes())

	expectedBlock := CodeBlock{
		Lang:   "sh",
		Script: "echo \"Hello World\"\n",
	}

	resultBlock, ok := doc.GetCodeBlock(expectedBlock.GenID())
	Diff(t, true, ok)
	Diff(t, expectedBlock, resultBlock)
}
