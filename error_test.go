package errors

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestA(t *testing.T) {
	_err := Annotate(sql.ErrNoRows,
		NotFound,
		Message("query user info"),
		DebugInfo{[]string{"app.invoke", "dao.query"}, "app a"},
		StackTrace("test a"))

	show(t, _err)

	t.Logf("Cause: %v", errors.Unwrap(_err))
}

func TestPredefined(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		err := NewNotFound("dao.QueryUserInfo",
			ResourceInfo{
				ResourceType: "user",
				ResourceName: "123",
				Owner:        "app",
				Description:  "user 123 not found",
			})

		err = fmt.Errorf("app.Handler: %w", err)
		err = fmt.Errorf("app.Dispatcher: %w", err)

		show(t, err)
	})

	t.Run("bad request", func(t *testing.T) {
		err := NewBadRequest("dao.UpdateItem",
			BadRequest{[]FieldViolation{
				{"username", "is required"},
				{"email", "is invalid"},
			}},
		)
		show(t, err)
	})

	t.Run("failed precondition", func(t *testing.T) {
		err := NewFailedPrecondition("app.any",
			PreconditionFailure{[]TypedViolation{
				{"TOC", "redis03", "no memory"},
				{"TOC", "redis03", "no cpu"},
			}},
		)
		show(t, err)
	})

	t.Run("failed precondition", func(t *testing.T) {
		err := NewUnauthenticated("app.any",
			ErrorInfo{
				Reason: "access_token expired",
				Domain: "taobao",
				Metadata: map[string]string{
					"nick": "user1",
					"age":  "10",
				},
			},
		)
		show(t, err)
	})

	t.Run("resource exhausted", func(t *testing.T) {
		err := NewResourceExhausted("app.any",
			QuotaFailure{[]Violation{{"rds", "concurrent max limit 15"}}})
		show(t, err)
	})

	t.Run("resource exhausted", func(t *testing.T) {
		err := NewResourceExhausted("app.any",
			QuotaFailure{[]Violation{{"rds", "concurrent max limit 15"}}})

		err = Annotate(err,
			RequestInfo{"abd123", ""},
			LocalizedMessage{"zh-CN", "中文提示"},
			Help{[]Link{{"about", "https://blog.igota.net/about"}}},
		)

		show(t, err)
	})
}

func TestDecode(t *testing.T) {
	const raw = `{"error":{"code":404,"message":"query user info: sql: no rows in result set","status":"NOT_FOUND","details":[{"@type":"type.googleapis.com/google.rpc.DebugInfo","stackEntries":["app.invoke","dao.query"],"detail":"app a"},{"@type":"type.googleapis.com/google.rpc.DebugInfo","stackEntries":["goroutine 6 [running]:","runtime/debug.Stack(0x1, 0x1, 0x1)","\tC:/code/go/src/runtime/debug/stack.go:24 +0xa5","github.com/gota33/errors.StackTrace.Annotate(0x10d48f6, 0x6, 0x111bd50, 0xc00005e840)","\tC:/workspace/github/gota33/errors/detail.go:365 +0x2d","github.com/gota33/errors.Annotate(0x11199a0, 0xc000050c50, 0xc00005df30, 0x4, 0x4, 0x1205ee0, 0xfb8a46)","\tC:/workspace/github/gota33/errors/errors.go:77 +0xbc","github.com/gota33/errors.TestA(0xc000045080)","\tC:/workspace/github/gota33/errors/error_test.go:12 +0x19e","testing.tRunner(0xc000045080, 0x10e9350)","\tC:/code/go/src/testing/testing.go:1194 +0xef","created by testing.(*T).Run","\tC:/code/go/src/testing/testing.go:1239 +0x2b3",""],"detail":"test a"}]}}`

	err := Decode(strings.NewReader(raw))
	show(t, err)
}

func show(t *testing.T, _err error) {
	t.Helper()

	t.Logf("code: %v", Code(_err))
	t.Logf("details: %d", len(Details(_err)))
	t.Logf("short: %v", _err)
	t.Logf("full message:\n%+v", _err)

	var buf bytes.Buffer
	if err := Encode(&buf, _err); err != nil {
		t.Fatal(err)
	}
	t.Log(buf.String())
}
