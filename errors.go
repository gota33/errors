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
	AppendMessage(msg string)
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

func (e *annotated) AppendMessage(msg string) {
	if e.message == "" {
		e.message = msg
	} else {
		e.message = msg + ": " + e.message
	}
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
	err, ok := cause.(*annotated)
	if !ok {
		err = &annotated{cause: cause}
	}

	Message(cause.Error()).Annotate(err)

	for _, annotation := range annotations {
		annotation.Annotate(err)
	}
	return err
}

type Message string

func (a Message) Annotate(m Modifier) {
	m.AppendMessage(string(a))
}

func Code(err error) StatusCode {
	if err == nil {
		return OK
	}

	var a *annotated
	if errors.As(err, &a) {
		return a.code
	}
	return Unknown
}

func Details(err error) (out []Any) {
	travel(err, visitor{
		OnDetails: func(details []Any) {
			out = append(out, details...)
		},
	})
	return
}

func Flatten(err error) error {
	var o annotated

	travel(err, visitor{
		OnCode: func(code StatusCode) {
			if o.code == OK {
				o.code = code
			}
		},
		OnDetails: func(details []Any) {
			o.details = append(o.details, details...)
		},
		OnError: func(cur error) {
			if o.message == "" {
				o.message = cur.Error()
			}
			o.cause = cur
		},
	})

	if o.code == OK {
		switch o.cause {
		case context.Canceled:
			o.code = Cancelled
		case context.DeadlineExceeded:
			o.code = DeadlineExceeded
		}
	}
	return &o
}

type visitor struct {
	OnCode    func(code StatusCode)
	OnDetails func(details []Any)
	OnError   func(err error)
}

func travel(root error, v visitor) {
	cur := root
	for cur != nil {
		if a, ok := cur.(*annotated); ok {
			if v.OnCode != nil {
				v.OnCode(a.code)
			}
			if v.OnDetails != nil {
				v.OnDetails(a.details)
			}
		}
		if v.OnError != nil {
			v.OnError(cur)
		}
		cur = errors.Unwrap(cur)
	}
}
