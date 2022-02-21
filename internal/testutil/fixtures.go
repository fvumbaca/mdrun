package testutil

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const FixtureDIR = "fixtures"

func FixtureReader(t *testing.T, fixtureName string) io.ReadCloser {
	f, err := os.Open(filepath.Join(FixtureDIR, fixtureName))
	if err != nil {
		t.Error(err)
	}
	return f
}

func FixtureBytes(t *testing.T, fixtureName string) []byte {
	b, err := ioutil.ReadFile(filepath.Join(FixtureDIR, fixtureName))
	if err != nil {
		t.Error(err)
	}
	return b
}
