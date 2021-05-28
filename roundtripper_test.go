package errors

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundTripper(t *testing.T) {
	const response = `{"data": "response data"}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch strings.TrimPrefix(r.URL.Path, "/") {
		case "happy":
			_, _ = io.WriteString(w, response)
		case "sad":
			w.WriteHeader(fullError.code.Http())
			_, _ = io.WriteString(w, rawFullError)
		case "internal":
			w.WriteHeader(fullError.code.Http())
			_, _ = io.WriteString(w, "123")
		case "plaintext":
			http.Error(w, sql.ErrNoRows.Error(), http.StatusNotFound)
		}
	}))
	defer srv.Close()

	client := http.Client{
		Transport: &RoundTripper{
			Parent: http.DefaultTransport,
		},
	}

	t.Run("happy", func(t *testing.T) {
		resp, err := client.Get(srv.URL + "/happy")
		if !assert.NoError(t, err) {
			return
		}

		data, err := io.ReadAll(resp.Body)
		if !assert.NoError(t, err) {
			return
		}

		assert.Equal(t, response, string(data))
	})

	t.Run("sad", func(t *testing.T) {
		_, err := client.Get(srv.URL + "/sad")
		if !assert.Error(t, err) {
			return
		}

		var a *annotated
		if assert.True(t, errors.As(err, &a)) {
			messageEx := fmt.Sprintf("Get \"%s/sad\": %s", srv.URL, fullError.Error())

			assert.Equal(t, fullError.code, Code(err))
			assert.ErrorIs(t, a.cause, fullError.code)
			assert.Equal(t, messageEx, err.Error())
		}
	})

	t.Run("internal", func(t *testing.T) {
		_, err := client.Get(srv.URL + "/internal")
		if !assert.Error(t, err) {
			return
		}

		var jErr *json.UnmarshalTypeError
		if !assert.ErrorAs(t, err, &jErr) {
			return
		}

		uErr := &url.Error{Op: "Get", URL: srv.URL + "/internal", Err: jErr}
		assert.Equal(t, Unknown, Code(err))
		assert.Equal(t, uErr.Error(), err.Error())
	})

	t.Run("plaintext", func(t *testing.T) {
		_, err := client.Get(srv.URL + "/plaintext")
		if !assert.Error(t, err) {
			return
		}

		cause := &url.Error{
			Op:  "Get",
			URL: srv.URL + "/plaintext",
			Err: fmt.Errorf("%s: %w", sql.ErrNoRows, Unknown),
		}
		assert.ErrorIs(t, err, Unknown)
		assert.Equal(t, Unknown, Code(err))
		assert.Equal(t, cause.Error(), err.Error())
	})
}
