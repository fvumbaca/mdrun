package testutil

import "testing"

func NoErr(t *testing.T, err error) {
	if err != nil {
		t.Error("unexpected error:", err)
	}
}
