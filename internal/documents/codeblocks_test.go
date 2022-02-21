package documents

import (
	"testing"

	"github.com/fvumbaca/mdrun/internal/testutil"
)

func TestCodeBlockGenID(t *testing.T) {
	block := CodeBlock{
		Lang:   "example",
		Script: "example",
	}
	testutil.Diff(t, block.GenID(), "ZXhhbXBs")
}
