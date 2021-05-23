package errors

import (
	"fmt"
	"net/http"
	"strconv"
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

func StrToCode(status string) StatusCode {
	for i, value := range statusList {
		if status == value {
			return StatusCode(i)
		}
	}
	return Unknown
}
