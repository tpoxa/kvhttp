package store

import (
  "context"
  "errors"
)

var ErrCheckSumInvalid = errors.New("checksum invalid")

//go:generate mockery --name=IStore --filename=store.go
type IStore interface {
  Get(ctx context.Context, key string) ([]byte, bool)
  Set(ctx context.Context, key string, value []byte)
  SetIfCheckSumMatch(ctx context.Context, key string, value []byte, checksum string) error
  Delete(ctx context.Context, key string)
}
