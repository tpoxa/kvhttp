package hasher

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
)

type Sha256 struct {
}

func NewSha256() *Sha256 {
	return &Sha256{}
}

func (s Sha256) Hash(_ context.Context, in []byte) (string, error) {
	h := sha256.New()
	_, err := h.Write(in)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
