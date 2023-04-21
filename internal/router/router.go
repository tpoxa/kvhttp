package router

import (
  "fmt"
  "io"
  "kvhttp/internal/hasher"
  "kvhttp/internal/store"
  "log"
  "net/http"
  "strings"
)

const maxLength = 1 << 3 //1mb

var (
  errInvalidKey    = fmt.Errorf("invalid key")
  errNotExists     = fmt.Errorf("not found")
  errInvalidLength = fmt.Errorf("invalid length")
)

type Router struct {
  store  store.IStore
  hasher hasher.Hasher
}

// NewRouter creates new instance
func NewRouter(store store.IStore, h hasher.Hasher) *Router {
  return &Router{store: store, hasher: h}
}

// ServeHTTP Serves requests implements http.Handler interface
// @todo health check?
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  log.Printf("%s %s", req.Method, req.URL.Path)
  ctx := req.Context()

  switch req.Method {
  case http.MethodPost, http.MethodPut, http.MethodPatch:
    k, err := getKey(req.URL.Path)
    if err != nil {
      handleError(w, err)
      return
    }
    if req.ContentLength == 0 || req.ContentLength > maxLength {
      handleError(w, errInvalidLength)
      return
    }

    body, err := io.ReadAll(io.LimitReader(req.Body, req.ContentLength))
    defer req.Body.Close()

    if err != nil {
      handleError(w, err)
      return
    }

    checkSum := req.Header.Get("If-Match")
    if checkSum != "" {
      if err = r.store.SetIfCheckSumMatch(ctx, k, body, checkSum); err != nil {
        handleError(w, err)
      }
    } else {
      r.store.Set(ctx, k, body)
    }

    checkSum, err = r.hasher.Hash(ctx, body)
    if err != nil {
      handleError(w, err)
      return
    }

    w.Header().Set("ETag", checkSum)

  case http.MethodGet:
    k, err := getKey(req.URL.Path)
    if err != nil {
      handleError(w, err)
      return
    }
    v, exists := r.store.Get(ctx, k)
    if !exists {
      handleError(w, errNotExists)
      return
    }

    checkSum, err := r.hasher.Hash(ctx, v)
    if err != nil {
      handleError(w, errNotExists)
      return
    }

    w.Header().Set("ETag", checkSum)
    w.Write(v)

  case http.MethodHead:
    k, err := getKey(req.URL.Path)
    if err != nil {
      handleError(w, err)
      return
    }
    v, exists := r.store.Get(ctx, k)
    if !exists {
      handleError(w, errNotExists)
      return
    }

    checkSum, err := r.hasher.Hash(ctx, v)
    if err != nil {
      handleError(w, errNotExists)
      return
    }

    w.Header().Set("ETag", checkSum)

    //@todo add e-tag support
  case http.MethodDelete:
    k, err := getKey(req.URL.Path)
    if err != nil {
      handleError(w, err)
      return
    }
    r.store.Delete(ctx, k)
    w.WriteHeader(http.StatusOK)
  default:
    http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
  }
}

// getKey gets store key from request's URL, may return errInvalidKey error
func getKey(path string) (string, error) {
  key := strings.Trim(path, "/")
  if key == "" {
    return "", errInvalidKey
  }
  return key, nil
}

// handleError determine error and write http response
func handleError(w http.ResponseWriter, err error) {
  switch err {
  case errInvalidKey:
    http.Error(w, "invalid key", http.StatusBadRequest)
  case errNotExists:
    http.Error(w, "not found", http.StatusNotFound)
  case errInvalidLength:
    http.Error(w, "invalid length", http.StatusBadRequest)
  case store.ErrCheckSumInvalid:
    http.Error(w, "e-tag mismatch", http.StatusConflict)
  default:
    http.Error(w, fmt.Sprintf("internal error"), http.StatusInternalServerError)
  }
}
