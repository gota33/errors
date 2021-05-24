package errors

import (
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
	details []Any
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

	WithMessage(cause.Error())(err)

	for _, annotation := range annotations {
		annotation(err)
	}
	return err
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

func WithDetail(details ...Any) Annotation {
	return func(err *annotated) { err.details = append(err.details, details...) }
}

func WithStack() Annotation {
	entries := strings.Split(string(debug.Stack()), "\n")
	detail := DebugInfo{StackEntries: entries}
	return WithDetail(detail)
}

func WithRequestInfo(requestId, servingData string) Annotation {
	return WithDetail(RequestInfo{requestId, servingData})
}

func WithLocalizedMessage(local string, message string) Annotation {
	return WithDetail(LocalizedMessage{local, message})
}

func WithHelp(links ...Link) Annotation {
	return WithDetail(Help{links})
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

func Details(err error) []Any {
	if err == nil {
		return nil
	}

	err = Flatten(err)

	var a *annotated
	if errors.As(err, &a) {
		return a.details
	}
	return nil
}

func Flatten(err error) error {
	var (
		o   annotated
		cur = err
	)
	for cur != nil {
		if a, ok := cur.(*annotated); ok {
			if o.code == 0 {
				o.code = a.code
			}
			o.details = append(o.details, a.details...)
		}
		if o.message == "" {
			o.message = cur.Error()
		}
		o.cause = cur
		cur = errors.Unwrap(cur)
	}
	if o.cause == nil {
		return nil
	}

	return &o
}
