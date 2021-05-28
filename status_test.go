package errors

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusCode(t *testing.T) {
	t.Run("annotate", func(t *testing.T) {
		for i := 0; i < int(totalStatus); i++ {
			s := StatusCode(i)
			m := &MockModifier{}
			s.Annotate(m)
			assert.Equal(t, s, m.Code)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		s := StatusCode(-1)
		assert.False(t, s.Valid())
	})

	t.Run("message", func(t *testing.T) {
		for i := 0; i < int(totalStatus); i++ {
			s := StatusCode(i)
			msg := fmt.Sprintf("%d %s", httpList[i], statusList[i])
			assert.Equal(t, msg, s.Error())
		}

		s := StatusCode(-1)
		assert.Equal(t, "500 Code(-1)", s.Error())
	})
}

func TestStatusName(t *testing.T) {
	t.Run("code", func(t *testing.T) {
		for i := 0; i < int(totalStatus); i++ {
			code := statusList[i].StatusCode()
			assert.Equal(t, StatusCode(i), code)
		}

		name := StatusName("don't known")
		assert.Equal(t, Unknown, name.StatusCode())
	})

	t.Run("name", func(t *testing.T) {
		for i := 0; i < int(totalStatus); i++ {
			lower := strings.ToLower(string(statusList[i]))
			code := StatusName(lower).StatusCode()
			assert.Equal(t, StatusCode(i), code)
		}
	})
}

type MockModifier struct {
	Code     StatusCode
	Messages []string
	Details  []Any
}

func (m *MockModifier) SetCode(code StatusCode) {
	m.Code = code
}

func (m *MockModifier) WrapMessage(msg string) {
	m.Messages = append(m.Messages, msg)
}

func (m *MockModifier) AppendDetails(details ...Any) {
	m.Details = append(m.Details, details...)
}
