package rundoc

import (
	"flag"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestParseCodeBloc_None(t *testing.T) {
	doc, err := Parse([]byte("no blocks in here"))
	noErr(t, err)
	diff(t, 0, len(doc.blocks))
}

func TestParseCodeBloc_Some(t *testing.T) {
	doc, err := Parse(loadFixture(t, "two_code_blocks.md"))
	noErr(t, err)
	diff(t, 2, len(doc.blocks))
}
