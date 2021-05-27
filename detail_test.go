package errors

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const typeUrlCustom = "custom/type"

var (
	resourceInfo        = ResourceInfo{"1", "2", "3", "4"}
	badRequest          = BadRequest{[]FieldViolation{{"1", "2"}, {"3", "4"}}}
	preconditionFailure = PreconditionFailure{[]TypedViolation{{"1", "2", "3"}, {"4", "5", "6"}}}
	errInfo             = ErrorInfo{"1", "2", map[string]string{"3": "4"}}
	quotaFailure        = QuotaFailure{[]Violation{{"1", "2"}, {"3", "4"}}}
	debugInfo           = DebugInfo{[]string{"1", "2"}, "3"}
	requestInfo         = RequestInfo{"1", "2"}
	help                = Help{[]Link{{"1", "2"}, {"3", "4"}}}
	localizedMessage    = LocalizedMessage{"1", "2"}
	any                 = AnyDetail{"@type": typeUrlCustom, "1": "2"}
)

type Detail struct {
	Any
	Type string
	Json string
}

var details = []Detail{
	{debugInfo, TypeUrlDebugInfo, `{"@type":"type.googleapis.com/google.rpc.DebugInfo","stackEntries":["1","2"],"detail":"3"}`},
	{resourceInfo, TypeUrlResourceInfo, `{"@type":"type.googleapis.com/google.rpc.ResourceInfo","resourceType":"1","resourceName":"2","owner":"3","description":"4"}`},
	{badRequest, TypeUrlBadRequest, `{"@type":"type.googleapis.com/google.rpc.BadRequest","fieldViolations":[{"field":"1","description":"2"},{"field":"3","description":"4"}]}`},
	{preconditionFailure, TypeUrlPreconditionFailure, `{"@type":"type.googleapis.com/google.rpc.PreconditionFailure","violations":[{"type":"1","subject":"2","description":"3"},{"type":"4","subject":"5","description":"6"}]}`},
	{errInfo, TypeUrlErrorInfo, `{"@type":"type.googleapis.com/google.rpc.ErrorInfo","reason":"1","domain":"2","metadata":{"3":"4"}}`},
	{quotaFailure, TypeUrlQuotaFailure, `{"@type":"type.googleapis.com/google.rpc.QuotaFailure","violations":[{"subject":"1","description":"2"},{"subject":"3","description":"4"}]}`},
	{requestInfo, TypeUrlRequestInfo, `{"@type":"type.googleapis.com/google.rpc.RequestInfo","requestId":"1","servingData":"2"}`},
	{help, TypeUrlHelp, `{"@type":"type.googleapis.com/google.rpc.Help","links":[{"description":"1","url":"2"},{"description":"3","url":"4"}]}`},
	{localizedMessage, TypeUrlLocalizedMessage, `{"@type":"type.googleapis.com/google.rpc.LocalizedMessage","local":"1","message":"2"}`},
	{any, typeUrlCustom, `{"@type":"custom/type","1":"2"}`},
}

func TestDetail(t *testing.T) {
	t.Run("type", func(t *testing.T) {
		for _, detail := range details {
			assert.Equal(t, detail.Type, detail.TypeUrl())
		}
	})

	t.Run("annotate", func(t *testing.T) {
		for _, detail := range details {
			m := &MockModifier{}
			detail.Annotate(m)
			assert.Equal(t, m.Details, []Any{detail.Any})
		}
	})

	t.Run("json", func(t *testing.T) {
		for _, detail := range details {
			data, err := json.Marshal(detail.Any)
			if assert.NoError(t, err) {
				assert.JSONEq(t, detail.Json, string(data))
			}
		}
	})

	t.Run("format", func(t *testing.T) {
		for _, detail := range details {
			strs := []string{
				fmt.Sprintf("%+v", detail),
				fmt.Sprintf("%v", detail),
				fmt.Sprintf("%s", detail),
				fmt.Sprintf("%q", detail),
			}
			for _, str := range strs {
				assert.NotEmpty(t, str)
				assert.NotContains(t, str, "!")
			}
		}
	})
}

func TestStackTrace(t *testing.T) {
	m := &MockModifier{}
	stack := StackTrace(msg)
	stack.Annotate(m)

	assert.NotEmpty(t, m.Details)

	detail := m.Details[0].(DebugInfo)
	assert.Equal(t, msg, detail.Detail)
	assert.NotEmpty(t, detail.StackEntries)

}
