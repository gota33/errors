package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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
			messageEx := fmt.Sprintf("GET %s/sad: %s", srv.URL, fullError.Error())

			assert.Equal(t, fullError.code, a.code)
			assert.Equal(t, fullError.code, a.cause)
			assert.Equal(t, messageEx, a.Error())
		}
	})

	t.Run("internal", func(t *testing.T) {
		_, err := client.Get(srv.URL + "/internal")
		if !assert.Error(t, err) {
			return
		}

		var a *annotated
		if assert.True(t, errors.As(err, &a)) {
			var jErr *json.UnmarshalTypeError
			if !assert.True(t, errors.As(err, &jErr)) {
				return
			}

			messageEx := fmt.Sprintf("GET %s/internal: %s", srv.URL, jErr)

			assert.Equal(t, Unknown, a.code)
			assert.True(t, errors.Is(a.cause, jErr))
			assert.Equal(t, messageEx, a.Error())
		}
	})
}
