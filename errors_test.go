package errors

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	t.Run("annotate nil", func(t *testing.T) {
		err := Annotate(nil, NotFound)
		assert.NoError(t, err)
	})

	code := NotFound
	cause := fmt.Errorf("wrapper: %w", sql.ErrNoRows)
	err := Annotate(cause, Message(msg), code, resourceInfo)
	a, ok := err.(*annotated)
	if !assert.True(t, ok) {
		return
	}

	t.Run("annotate", func(t *testing.T) {
		messageEx := fmt.Sprintf("%s: %s", msg, cause)

		assert.Equal(t, NotFound, a.code)
		assert.Equal(t, cause, a.cause)
		assert.Equal(t, messageEx, a.message)
	})

	t.Run("code", func(t *testing.T) {
		assert.Equal(t, code, Code(err))
		assert.Equal(t, OK, Code(nil))
		assert.Equal(t, Unknown, Code(errors.New("unknown error")))
	})

	t.Run("code2", func(t *testing.T) {
		_err := Annotate(NotFound, Message("user"))
		assert.Equal(t, NotFound, Code(_err))
	})

	t.Run("code3", func(t *testing.T) {
		assert.Equal(t, Cancelled, Code(context.Canceled))
		assert.Equal(t, DeadlineExceeded, Code(context.DeadlineExceeded))
	})

	t.Run("detail", func(t *testing.T) {
		assert.Equal(t, []Any{resourceInfo}, Details(err))
	})

	t.Run("format", func(t *testing.T) {
		strs := []string{
			fmt.Sprintf("%v", err),
			fmt.Sprintf("%+v", err),
			fmt.Sprintf("%q", err),
		}

		for _, str := range strs {
			assert.NotEmpty(t, str)
			assert.NotContains(t, str, "!")
		}
	})
}

func TestTemporary(t *testing.T) {
	items := map[error]bool{
		OK:                                  false,
		Internal:                            false,
		fmt.Errorf("wrap: %w", Internal):    false,
		Unavailable:                         true,
		fmt.Errorf("wrap: %w", Unavailable): true,
		Annotate(&net.DNSError{IsTemporary: true}, Internal): true,
	}
	for err, ok := range items {
		assert.Equal(t, ok, Temporary(err), err.Error())
	}
}

func TestFlatten(t *testing.T) {
	causes := map[StatusCode]error{
		Cancelled:        context.Canceled,
		DeadlineExceeded: context.DeadlineExceeded,
		Unknown:          errors.New("unknown error"),
		OK:               nil,
	}

	hide := func(typeUrl string) DetailMapper {
		return func(a Any) Any {
			if typeUrl != a.TypeUrl() {
				return a
			}
			return nil
		}
	}

	mappers := detailMappers{
		hide(TypeUrlDebugInfo),
		hide(TypeUrlResourceInfo),
	}

	for code, cause := range causes {
		err := Annotate(cause, Message(msg), debugInfo, resourceInfo, errInfo)
		err = Flatten(err, mappers...)

		a, ok := err.(*annotated)
		if !assert.True(t, ok) {
			return
		}

		assert.Equal(t, code, a.code)

		for _, detail := range details {
			assert.NotEqual(t, debugInfo, detail)
			assert.NotEqual(t, resourceInfo, detail)
		}
	}
}
