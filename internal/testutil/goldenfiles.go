package testutil

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var updateGoldenFiles = false

func init() {
	if _, set := os.LookupEnv("GOLDEN_UPDATE"); set {
		updateGoldenFiles = true
	}
}

// GoldenFileDiff will compare the provided result with a named golden file.
// When required, all golden files can be updated by setting `GOLDEN_UPDATE=1`
// before running tests using golden files.
func GoldenFileDiff(t *testing.T, filename string, result []byte) {
	if updateGoldenFiles {
		err := ioutil.WriteFile(filepath.Join(FixtureDIR, filename), result, 0664)
		NoErr(t, err)
	}
	b, err := ioutil.ReadFile(filepath.Join(FixtureDIR, filename))
	NoErr(t, err)
	Diff(t, b, result)
}
