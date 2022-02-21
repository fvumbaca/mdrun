package documents

import (
	"crypto/sha256"
	"encoding/base64"
)

type CodeBlock struct {
	Lang   string
	Script string
}

func (cb CodeBlock) GenID() string {
	// TODO: Review this to minimize collisions
	id := base64.RawStdEncoding.EncodeToString(sha256.New().Sum([]byte(cb.Lang +
		string(cb.Script))))
	if len(id) < 8 {
		return id
	}
	return id[:8]
}
