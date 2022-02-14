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
	{"basic.md", "basic.html.golden"},
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestBasicBuildHTML(t *testing.T) {
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
				ioutil.WriteFile(golden, buff.Bytes(), 0655)
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

func loadFixture(t *testing.T, name string) []byte {
	content, err := ioutil.ReadFile("fixtures/" + name)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	return content
}
