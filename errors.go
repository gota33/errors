package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"runtime/debug"
	"strings"
)

type annotated struct {
	cause   error
	code    StatusCode
	message string
	details []Detail
}

func (e annotated) Unwrap() error { return e.cause }

func (e annotated) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Code    int      `json:"code"`
		Message string   `json:"message,omitempty"`
		Status  string   `json:"status"`
		Details []Detail `json:"details,omitempty"`
	}{
		Code:    e.code.HttpCode(),
		Message: e.Error(),
		Status:  e.code.String(),
		Details: e.details,
	})
}

func (e annotated) Error() string {
	if e.message == "" {
		return e.cause.Error()
	}
	return e.message + ": " + e.cause.Error()
}

func (e annotated) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			_, _ = fmt.Fprintf(f, "status: %q\n", e.code.Error())
			_, _ = fmt.Fprintf(f, "message: %q\n", e.Error())
			for i, detail := range e.details {
				_, _ = fmt.Fprintf(f, "detail[%d]:\n%+v\n", i, detail)
			}
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(f, e.Error())
	case 'q':
		_, _ = fmt.Fprintf(f, "%q", e.Error())
	}
}

type Annotation func(err *annotated)

func WithCode(code StatusCode) Annotation {
	return func(err *annotated) { err.code = code }
}

func WithMessage(message string) Annotation {
	return func(err *annotated) {
		if err.message == "" {
			err.message = message
		} else {
			err.message = message + ": " + err.message
		}
	}
}

func WithDetail(details ...Detail) Annotation {
	return func(err *annotated) { err.details = append(err.details, details...) }
}

func WithStack() Annotation {
	return func(err *annotated) {
		entries := strings.Split(string(debug.Stack()), "\n")
		detail := DebugInfo{StackEntries: entries}
		err.details = append(err.details, detail)
	}
}

func WithRequestInfo(requestId, servingData string) Annotation {
	return func(err *annotated) {
		err.details = append(err.details, RequestInfo{requestId, servingData})
	}
}

func WithLocalizedMessage(local string, message string) Annotation {
	return func(err *annotated) {
		err.details = append(err.details, LocalizedMessage{local, message})
	}
}

func WithHelp(links ...Link) Annotation {
	return func(err *annotated) {
		err.details = append(err.details, Help{links})
	}
}

func Annotate(cause error, annotations ...Annotation) error {
	err, ok := cause.(*annotated)
	if !ok {
		err = &annotated{cause: cause}
	}
	for _, annotation := range annotations {
		annotation(err)
	}
	return err
}

func Details(err error) (details []Detail) {
	if a, ok := err.(*annotated); ok {
		return a.details
	}
	return
}

type encoded struct {
	Code    int             `json:"code"`
	Message string          `json:"message,omitempty"`
	Status  string          `json:"status"`
	Details json.RawMessage `json:"details,omitempty"`
}

func Encode(err error) (data []byte, _err error) {
	var o struct {
		Code    int      `json:"code"`
		Message string   `json:"message,omitempty"`
		Status  string   `json:"status"`
		Details []Detail `json:"details,omitempty"`
	}

	cur := err
	for cur != nil {
		if a, ok := cur.(*annotated); ok {
			if o.Code == 0 {
				o.Code = a.code.HttpCode()
				o.Status = a.code.String()
			}
			o.Details = append(o.Details, a.details...)
		}
		if o.Message != "" {
			o.Message += ": "
		}
		o.Message += cur.Error()
		cur = errors.Unwrap(cur)
	}

	return json.Marshal(o)
}

func Decode(body io.Reader) (err error) {
	var msg encoded
	if err = json.NewDecoder(body).Decode(&msg); err != nil {
		return
	}

	details, err := decodeDetails(msg.Details)
	if err != nil {
		return
	}

	code := StrToCode(msg.Status)

	return &annotated{
		cause:   code,
		code:    code,
		message: msg.Message,
		details: details,
	}
}
