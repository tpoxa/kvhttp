package router

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	hashermocks "kvhttp/internal/hasher/mocks"
	"kvhttp/internal/store/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouter_ServeHTTP(t *testing.T) {

	type fields struct {
		store  *mocks.IStore
		hasher *hashermocks.Hasher
	}

	type args struct {
		w   *httptest.ResponseRecorder
		req *http.Request
	}

	type expected struct {
		code           int    // we expect certain http code
		verifyResponse bool   // shall we check response data
		resp           []byte // we expect response
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		do       func(fields)
		expected *expected
	}{
		{
			name: "GET request should trigger store get, record not found returns 404",
			fields: fields{
				store:  mocks.NewIStore(t),
				hasher: hashermocks.NewHasher(t),
			},
			args: args{
				req: httptest.NewRequest(http.MethodGet, "/foo", nil),
				w:   httptest.NewRecorder(),
			},
			expected: &expected{
				code: http.StatusNotFound,
			},
			do: func(fields fields) {
				fields.store.On("Get", mock.Anything, "foo").Return(nil, false).Times(1) // exactly 1 time
			},
		},
		{
			name: "GET request should trigger store get, record not found returns exact same value as store, statusOK",
			fields: fields{
				store:  mocks.NewIStore(t),
				hasher: hashermocks.NewHasher(t),
			},
			args: args{
				req: httptest.NewRequest(http.MethodGet, "/foo", nil),
				w:   httptest.NewRecorder(),
			},
			expected: &expected{
				code:           http.StatusOK,
				verifyResponse: true,
				resp:           []byte("bar"),
			},
			do: func(fields fields) {
				fields.store.On("Get", mock.Anything, "foo").Return([]byte("bar"), true).Times(1) // exactly 1 time
				fields.hasher.On("Hash", mock.Anything, []byte("bar")).Return("fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9", nil)
			},
		},
		{
			name: "POST request should trigger store get",
			fields: fields{
				store:  mocks.NewIStore(t),
				hasher: hashermocks.NewHasher(t),
			},
			args: args{
				req: httptest.NewRequest(http.MethodPost, "/foo", strings.NewReader("bar")),
				w:   httptest.NewRecorder(),
			},
			expected: &expected{
				code: http.StatusOK, // set should go without error
			},
			do: func(fields fields) {
				fields.store.On("Set", mock.Anything, "foo", []byte("bar")).Return().Times(1) // exactly 1 time
				fields.hasher.On("Hash", mock.Anything, []byte("bar")).Return("fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9", nil)

			},
		},
		{
			name: "PUT request should trigger store get",
			fields: fields{
				store:  mocks.NewIStore(t),
				hasher: hashermocks.NewHasher(t),
			},
			args: args{
				req: httptest.NewRequest(http.MethodPut, "/foo", strings.NewReader("bar")),
				w:   httptest.NewRecorder(),
			},
			expected: &expected{
				code: http.StatusOK, // set should go without error
			},
			do: func(fields fields) {
				fields.store.On("Set", mock.Anything, "foo", []byte("bar")).Return().Times(1) // exactly 1 time
				fields.hasher.On("Hash", mock.Anything, []byte("bar")).Return("fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9", nil)

			},
		},
		{
			name: "PATCH request should trigger store get",
			fields: fields{
				store:  mocks.NewIStore(t),
				hasher: hashermocks.NewHasher(t),
			},
			args: args{
				req: httptest.NewRequest(http.MethodPatch, "/foo", strings.NewReader("bar")),
				w:   httptest.NewRecorder(),
			},
			expected: &expected{
				code: http.StatusOK, // set should go without error
			},
			do: func(fields fields) {
				fields.store.On("Set", mock.Anything, "foo", []byte("bar")).Return().Times(1) // exactly 1 time
				fields.hasher.On("Hash", mock.Anything, []byte("bar")).Return("fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9", nil)

			},
		},
		{
			name: "DELETE request should trigger store delete",
			fields: fields{
				store:  mocks.NewIStore(t),
				hasher: hashermocks.NewHasher(t),
			},
			args: args{
				req: httptest.NewRequest(http.MethodDelete, "/foo", nil),
				w:   httptest.NewRecorder(),
			},
			expected: &expected{
				code: http.StatusOK, // delete should go without error
			},
			do: func(fields fields) {
				fields.store.On("Delete", mock.Anything, "foo").Return().Times(1) // exactly 1 time
			},
		},
		{
			name: "Get request with empty URL path should cause bad request",
			fields: fields{
				store:  mocks.NewIStore(t),
				hasher: hashermocks.NewHasher(t),
			},
			args: args{
				req: httptest.NewRequest(http.MethodGet, "/", nil),
				w:   httptest.NewRecorder(),
			},
			do: func(fields fields) {

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.do(tt.fields)

			r := NewRouter(tt.fields.store, tt.fields.hasher)
			r.ServeHTTP(tt.args.w, tt.args.req)
			if tt.expected != nil {
				resp := tt.args.w.Result()
				assert.Equalf(t, tt.expected.code, resp.StatusCode, "HTTP status code mismatch")

				if tt.expected.verifyResponse {
					data, err := io.ReadAll(resp.Body)
					if err != nil {
						t.Errorf("expected error to be nil got %v", err)
					}
					assert.Equalf(t, tt.expected.resp, data, "expected response is different")
				}
			}
		})
	}
}
