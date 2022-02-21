package testutil

import "testing"

func NoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Error("unexpected error:", err)
	}
}
