package mem

import (
  "context"
  "kvhttp/internal/hasher"
  "kvhttp/internal/store"
  "sync"
)

type Store struct {
  records map[string][]byte
  lock    *sync.RWMutex
  hasher  hasher.Hasher
}

func NewStore(h hasher.Hasher) *Store {
  return &Store{
    records: make(map[string][]byte),
    lock:    &sync.RWMutex{},
    hasher:  h,
  }
}

func (s *Store) Get(ctx context.Context, key string) ([]byte, bool) {
  s.lock.RLock()
  defer s.lock.RUnlock()
  value, ok := s.records[key]
  return value, ok
}

func (s *Store) Set(ctx context.Context, key string, value []byte) {
  s.lock.Lock()
  defer s.lock.Unlock()
  s.records[key] = value
}

func (s *Store) Delete(ctx context.Context, key string) {
  s.lock.Lock()
  defer s.lock.Unlock()
  delete(s.records, key)
}

func (s *Store) SetIfCheckSumMatch(ctx context.Context, key string, value []byte, checkSum string) error {
  s.lock.Lock()
  defer s.lock.Unlock()

  recordValue, ok := s.records[key]
  if !ok {
    return store.ErrCheckSumInvalid
  }
  recordCheckSum, err := s.hasher.Hash(ctx, recordValue)
  if err != nil {
    return err
  }
  if recordCheckSum != checkSum {
    return store.ErrCheckSumInvalid
  }

  s.records[key] = value
  return nil
}
