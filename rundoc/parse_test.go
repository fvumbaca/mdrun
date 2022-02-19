package rundoc

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var update = flag.Bool("update", false, "Update golden file for html building tests")

var tableRenderTests = [][]string{
	{"basic.md", "basic.golden.html"},
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestBasicBuildBodyHTML(t *testing.T) {
	for i, caseFiles := range tableRenderTests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			input := loadFixture(t, caseFiles[0])

			doc, err := Parse(input)
			if err != nil {
				t.Error(err)
			}

			var buff bytes.Buffer
			doc.WriteHTML(&buff)

			golden := filepath.Join("fixtures", caseFiles[1])
			if *update {
				ioutil.WriteFile(golden, buff.Bytes(), 0664)
			}

			expected, err := ioutil.ReadFile(golden)
			if err != nil {
				t.Error(err)
			}
			if diff := cmp.Diff(string(expected), string(buff.Bytes())); diff != "" {
				t.Error(diff)
			}
		})
	}
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
