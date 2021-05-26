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
)

const totalStatus = 17

var (
	statusList = [totalStatus]string{
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
	return fmt.Sprintf("%d %s", c.HttpCode(), c.String())
}

func (c StatusCode) String() string {
	if c.Valid() {
		return statusList[c]
	} else {
		return "Code(" + strconv.FormatInt(int64(c), 10) + ")"
	}
}

func (c StatusCode) HttpCode() int {
	if c.Valid() {
		return httpList[c]
	} else {
		return http.StatusInternalServerError
	}
}

type StatusCodeName string

func (s StatusCodeName) StatusCode() StatusCode {
	name := strings.ToUpper(string(s))
	for i, value := range statusList {
		if name == value {
			return StatusCode(i)
		}
	}
	return Unknown
}

type HttpStatusCode int

func (c HttpStatusCode) StatusCode() StatusCode {
	num := int(c)
	for i, value := range httpList {
		if num == value {
			return StatusCode(i)
		}
	}
	return Unknown
}

func (c HttpStatusCode) Annotate(m Modifier) {
	m.SetCode(c.StatusCode())
}
