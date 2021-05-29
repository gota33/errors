package errors

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	OK StatusCode = iota
	Cancelled
	Unknown
	InvalidArgument
	DeadlineExceeded
	NotFound
	AlreadyExists
	PermissionDenied
	ResourceExhausted
	FailedPrecondition
	Aborted
	OutOfRange
	Unimplemented
	Internal
	Unavailable
	DataLoss
	Unauthenticated

	totalStatus
)

var (
	statusList = [totalStatus]StatusName{
		"OK",
		"CANCELLED",
		"UNKNOWN",
		"INVALID_ARGUMENT",
		"DEADLINE_EXCEEDED",
		"NOT_FOUND",
		"ALREADY_EXISTS",
		"PERMISSION_DENIED",
		"RESOURCE_EXHAUSTED",
		"FAILED_PRECONDITION",
		"ABORTED",
		"OUT_OF_RANGE",
		"UNIMPLEMENTED",
		"INTERNAL",
		"UNAVAILABLE",
		"DATA_LOSS",
		"UNAUTHENTICATED",
	}
	httpList = [totalStatus]int{
		http.StatusOK,
		499,
		http.StatusInternalServerError,
		http.StatusBadRequest,
		http.StatusGatewayTimeout,
		http.StatusNotFound,
		http.StatusConflict,
		http.StatusForbidden,
		http.StatusTooManyRequests,
		http.StatusBadRequest,
		http.StatusConflict,
		http.StatusBadRequest,
		http.StatusNotImplemented,
		http.StatusInternalServerError,
		http.StatusServiceUnavailable,
		http.StatusInternalServerError,
		http.StatusUnauthorized,
	}
)

type StatusCode int

func (c StatusCode) Annotate(m Modifier) {
	m.SetCode(c)
}

func (c StatusCode) Valid() bool {
	return c >= 0 && c < totalStatus
}

func (c StatusCode) Error() string {
	return fmt.Sprintf("%d %s", c.Http(), c.String())
}

func (c StatusCode) Temporary() bool {
	return c == Unavailable
}

func (c StatusCode) Name() StatusName {
	if c.Valid() {
		return statusList[c]
	}
	return StatusName("Code(" + strconv.FormatInt(int64(c), 10) + ")")
}

func (c StatusCode) String() string { return string(c.Name()) }

func (c StatusCode) Http() int {
	if c.Valid() {
		return httpList[c]
	}
	return http.StatusInternalServerError
}

type StatusName string

func (s StatusName) StatusCode() StatusCode {
	name := StatusName(s.String())
	for i, value := range statusList {
		if name == value {
			return StatusCode(i)
		}
	}
	return Unknown
}

func (s StatusName) String() string {
	return strings.ToUpper(string(s))
}
