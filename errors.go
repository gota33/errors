package errors

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Modifier interface {
	SetCode(code StatusCode)
	WrapMessage(msg string)
	AppendDetails(details ...Any)
}

type Annotation interface {
	Annotate(m Modifier)
}

type annotated struct {
	cause   error
	code    StatusCode
	message string
	details []Any
}

func (e *annotated) SetCode(code StatusCode) {
	e.code = code
}

func (e *annotated) AppendDetails(details ...Any) {
	e.details = append(e.details, details...)
}

func (e *annotated) WrapMessage(msg string) {
	next := msg + ": " + e.message
	e.message = strings.TrimPrefix(next, ": ")
}

func (e annotated) Unwrap() error { return e.cause }
func (e annotated) Error() string { return e.message }

func (e annotated) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			_, _ = fmt.Fprintf(f, "status: %q\n", e.code.Error())
			_, _ = fmt.Fprintf(f, "message: %q\n", e.Error())
			for i, detail := range e.details {
				_, _ = fmt.Fprintf(f, "detail[%d]:\n", i)
				str := fmt.Sprintf("\t%+v", detail)
				str = strings.ReplaceAll(str, "\n", "\n\t")
				_, _ = io.WriteString(f, str[0:len(str)-1])
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

func Annotate(cause error, annotations ...Annotation) error {
	if cause == nil || len(annotations) == 0 {
		return cause
	}

	err := &annotated{
		cause:   cause,
		message: cause.Error(),
	}

	if code, ok := cause.(StatusCode); ok {
		err.code = code
	}

	for _, annotation := range annotations {
		annotation.Annotate(err)
	}
	return err
}

type Message string

func (a Message) Annotate(m Modifier) {
	m.WrapMessage(string(a))
}

func Code(err error) (out StatusCode) {
	travel(err, visitor{
		OnCode: func(code StatusCode) bool {
			out = code
			return out == OK
		},
	})

	if err == nil || out != OK {
		return
	}

	switch {
	case errors.Is(err, context.DeadlineExceeded):
		out = DeadlineExceeded
	case errors.Is(err, context.Canceled):
		out = Cancelled
	default:
		out = Unknown
	}
	return
}

func Temporary(err error) (out bool) {
	travel(err, visitor{
		OnError: func(err error) bool {
			if tErr, ok := err.(interface {
				Temporary() bool
			}); ok {
				out = tErr.Temporary()
				return false
			}
			return true
		},
	})
	return
}

func Details(err error) (out []Any) {
	travel(err, visitor{
		OnDetails: func(details []Any) bool {
			out = append(out, details...)
			return true
		},
	})
	return
}

func HideDebugInfo(a Any) Any {
	if a.TypeUrl() == TypeUrlDebugInfo {
		return nil
	}
	return a
}

type DetailMapper func(a Any) Any

type detailMappers []DetailMapper

func (fn detailMappers) Map(detail Any) (out Any) {
	out = detail
	for _, mapper := range fn {
		if out == nil {
			return
		}
		out = mapper(out)
	}
	return
}

func Flatten(err error, mappers ...DetailMapper) error {
	var o annotated

	travel(err, visitor{
		OnCode: func(code StatusCode) bool {
			if o.code == OK {
				o.code = code
			}
			return true
		},
		OnDetails: func(details []Any) bool {
			mapper := detailMappers(mappers)
			for _, detail := range details {
				if d := mapper.Map(detail); d != nil {
					o.details = append(o.details, d)
				}
			}
			return true
		},
		OnError: func(cur error) bool {
			if o.message == "" {
				o.message = cur.Error()
			}
			o.cause = cur
			return true
		},
	})

	if o.code == OK {
		switch o.cause {
		case nil:
			// Keep OK
		case context.Canceled:
			o.code = Cancelled
		case context.DeadlineExceeded:
			o.code = DeadlineExceeded
		default:
			o.code = Unknown
		}
	}
	return &o
}

type visitor struct {
	OnCode    func(code StatusCode) bool
	OnDetails func(details []Any) bool
	OnError   func(err error) bool
}

func travel(root error, v visitor) {
	cur := root
Loop:
	for cur != nil {
		switch a := cur.(type) {
		case StatusCode:
			if v.OnCode != nil && !v.OnCode(a) {
				break Loop
			}
		case *annotated:
			if v.OnCode != nil && !v.OnCode(a.code) {
				break Loop
			}
			if v.OnDetails != nil && !v.OnDetails(a.details) {
				break Loop
			}
		}
		if v.OnError != nil && !v.OnError(cur) {
			break Loop
		}
		cur = errors.Unwrap(cur)
	}
}
