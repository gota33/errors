package errors

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	rawFullError   = `{"error":{"code":500,"message":"msg: sql: connection is already closed","status":"INTERNAL","details":[{"@type":"type.googleapis.com/google.rpc.ResourceInfo","resourceType":"1","resourceName":"2","owner":"3","description":"4"},{"@type":"type.googleapis.com/google.rpc.BadRequest","fieldViolations":[{"field":"1","description":"2"},{"field":"3","description":"4"}]},{"@type":"type.googleapis.com/google.rpc.PreconditionFailure","violations":[{"type":"1","subject":"2","description":"3"},{"type":"4","subject":"5","description":"6"}]},{"@type":"type.googleapis.com/google.rpc.ErrorInfo","reason":"1","domain":"2","metadata":{"3":"4"}},{"@type":"type.googleapis.com/google.rpc.QuotaFailure","violations":[{"subject":"1","description":"2"},{"subject":"3","description":"4"}]},{"@type":"type.googleapis.com/google.rpc.DebugInfo","stackEntries":["1","2"],"detail":"3"},{"@type":"type.googleapis.com/google.rpc.RequestInfo","requestId":"1","servingData":"2"},{"@type":"type.googleapis.com/google.rpc.Help","links":[{"description":"1","url":"2"},{"description":"3","url":"4"}]},{"@type":"type.googleapis.com/google.rpc.LocalizedMessage","local":"1","message":"2"},{"1":"2","@type":"custom/type"}]}}`
	rawNoDebugInfo = `{"error":{"code":500,"message":"msg: sql: connection is already closed","status":"INTERNAL","details":[{"@type":"type.googleapis.com/google.rpc.ResourceInfo","resourceType":"1","resourceName":"2","owner":"3","description":"4"},{"@type":"type.googleapis.com/google.rpc.BadRequest","fieldViolations":[{"field":"1","description":"2"},{"field":"3","description":"4"}]},{"@type":"type.googleapis.com/google.rpc.PreconditionFailure","violations":[{"type":"1","subject":"2","description":"3"},{"type":"4","subject":"5","description":"6"}]},{"@type":"type.googleapis.com/google.rpc.ErrorInfo","reason":"1","domain":"2","metadata":{"3":"4"}},{"@type":"type.googleapis.com/google.rpc.QuotaFailure","violations":[{"subject":"1","description":"2"},{"subject":"3","description":"4"}]},{"@type":"type.googleapis.com/google.rpc.RequestInfo","requestId":"1","servingData":"2"},{"@type":"type.googleapis.com/google.rpc.Help","links":[{"description":"1","url":"2"},{"description":"3","url":"4"}]},{"@type":"type.googleapis.com/google.rpc.LocalizedMessage","local":"1","message":"2"},{"1":"2","@type":"custom/type"}]}}`
)

var fullError = &annotated{
	cause:   sql.ErrConnDone,
	code:    Internal,
	message: msg + ": " + sql.ErrConnDone.Error(),
	details: []Any{
		resourceInfo,
		badRequest,
		preconditionFailure,
		errInfo,
		quotaFailure,
		debugInfo,
		requestInfo,
		help,
		localizedMessage,
		any,
	},
}

func TestEncode(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		encode := func(raw string, filters ...DetailFilter) {
			var buf bytes.Buffer
			enc := NewEncoder(json.NewEncoder(&buf))
			enc.Filters = filters
			err := enc.Encode(fullError)

			if assert.NoError(t, err) {
				assert.Equal(t, raw, strings.TrimSpace(buf.String()))
			}
		}

		encode(rawFullError)
		encode(rawNoDebugInfo, HideDebugInfo)
	})

	t.Run("no encoder", func(t *testing.T) {
		enc := NewEncoder(nil)
		err := enc.Encode(Internal)
		assert.True(t, errors.Is(err, ErrNoEncoder))
	})
}

func TestDecode(t *testing.T) {
	t.Run("decode", func(t *testing.T) {
		r := strings.NewReader(rawFullError)
		dec := NewDecoder(json.NewDecoder(r))

		err := dec.Decode()
		assert.True(t, errors.Is(err, Internal))
		assert.True(t, strings.HasSuffix(err.Error(), fullError.cause.Error()))
		assert.Len(t, fullError.details, len(err.(*annotated).details))
	})

	t.Run("no decoder", func(t *testing.T) {
		dec := NewDecoder(nil)
		err := dec.Decode()
		assert.True(t, errors.Is(err, ErrNoDecoder))
	})
}

func TestCustomType(t *testing.T) {
	Register(typeUrlCustom, func() Any { return new(AnyDetail) })
	assert.Contains(t, typeProvider, typeUrlCustom)

	Register(typeUrlCustom, nil)
	assert.NotContains(t, typeProvider, typeUrlCustom)
}
