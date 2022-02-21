package testutil

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Diff(t *testing.T, expected, result interface{}) {
	if diff := cmp.Diff(expected, result); diff != "" {
		t.Error(diff)
	}
}

func DiffFromFixture(t *testing.T, fixtureName string, result []byte) {
	expected, err := ioutil.ReadFile(filepath.Join("fixtures", fixtureName))
	if err != nil {
		t.Error(err)
		return
	}
	Diff(t, expected, result)
}
