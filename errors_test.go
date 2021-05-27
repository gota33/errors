package errors

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
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

func TestFlatten(t *testing.T) {
	causes := map[StatusCode]error{
		Cancelled:        context.Canceled,
		DeadlineExceeded: context.DeadlineExceeded,
	}

	for code, cause := range causes {
		err := Annotate(cause, Message(msg), debugInfo)
		err = Flatten(err)

		a, ok := err.(*annotated)
		if !assert.True(t, ok) {
			return
		}

		assert.Equal(t, code, a.code)
	}
}
