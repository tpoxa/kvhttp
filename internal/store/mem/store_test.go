package mem

import (
	"context"
	"github.com/stretchr/testify/assert"
	hashermocks "kvhttp/internal/hasher/mocks"
	"testing"
)

func TestStore(t *testing.T) {
	ctx := context.Background()

	h := hashermocks.NewHasher(t)
	store := NewStore(h)

	_, exists := store.Get(ctx, "foo")
	if exists {
		t.Fatal("should not exist")
	}

	store.Set(ctx, "foo", []byte("bar"))

	val, exists := store.Get(ctx, "foo")
	if !exists {
		t.Fatal("should exist")
	}
	if !assert.Equal(t, []byte("bar"), val) {
		t.Fatal("value does not match")
	}

	store.Delete(ctx, "foo")
	_, exists = store.Get(ctx, "foo")

	if exists {
		t.Fatal("should not exist after deletion")
	}

}
