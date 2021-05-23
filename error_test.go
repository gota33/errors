package errors

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestA(t *testing.T) {
	cause := sql.ErrNoRows

	detail := DebugInfo{
		StackEntries: []string{"app.invoke", "dao.query"},
		Detail:       "something wrong",
	}

	_err := Annotate(cause,
		WithCode(NotFound),
		WithMessage("query user info"),
		WithDetail(detail),
		WithStack())

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
			FieldViolation{"username", "is required"},
			FieldViolation{"email", "is invalid"},
		)
		show(t, err)
	})

	t.Run("failed precondition", func(t *testing.T) {
		err := NewFailedPrecondition("app.any",
			TypedViolation{"TOC", "redis03", "no memory"},
			TypedViolation{"TOC", "redis03", "no cpu"},
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
			Violation{"rds", "concurrent max limit 15"})
		show(t, err)
	})

	t.Run("resource exhausted", func(t *testing.T) {
		err := NewResourceExhausted("app.any",
			Violation{"rds", "concurrent max limit 15"})

		err = Annotate(err,
			WithRequestInfo("abd123", ""),
			WithLocalizedMessage("zh-CN", "中文提示"),
			WithHelp(Link{"about", "https://blog.igota.net/about"}),
		)

		show(t, err)
	})
}

func TestDecode(t *testing.T) {
	const raw = `{"code":404,"message":"query user info: sql: no rows in result set","status":"NOT_FOUND","details":[{"@type":"type.googleapis.com/google.rpc.DebugInfo","stackEntries":["app.invoke","dao.query"],"detail":"something wrong"},{"@type":"type.googleapis.com/google.rpc.DebugInfo","stackEntries":["goroutine 6 [running]:","runtime/debug.Stack(0x1, 0xc00005de38, 0x107c03f)","\tC:/code/go/src/runtime/debug/stack.go:24 +0xa5","github.com/gota33/errors.WithStackTrace.func1(0xc00005e840)","\tC:/workspace/github/gota33/errors/errors.go:61 +0x3b","github.com/gota33/errors.Annotate(0x1115260, 0xc000050c50, 0xc00005df00, 0x4, 0x4, 0x112e7d4, 0xf)","\tC:/workspace/github/gota33/errors/errors.go:73 +0x62","github.com/gota33/errors.TestA(0xc000045080)","\tC:/workspace/github/gota33/errors/error_test.go:18 +0x251","testing.tRunner(0xc000045080, 0x10e7320)","\tC:/code/go/src/testing/testing.go:1194 +0xef","created by testing.(*T).Run","\tC:/code/go/src/testing/testing.go:1239 +0x2b3",""]}]}`

	err := Decode(strings.NewReader(raw))
	show(t, err)
}

func show(t *testing.T, _err error) {
	t.Helper()

	t.Logf("%v", _err)
	t.Logf("%+v", _err)

	data, err := Encode(_err)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))
}
