package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const msg = "msg"

var (
	resourceInfo        ResourceInfo
	badRequest          BadRequest
	preconditionFailure PreconditionFailure
	errInfo             ErrorInfo
	quotaFailure        QuotaFailure
	debugInfo           DebugInfo
)

var preErrors = []*annotated{
	{NotFound, NotFound, msg, []Any{resourceInfo}},
	{InvalidArgument, InvalidArgument, msg, []Any{badRequest}},
	{FailedPrecondition, FailedPrecondition, msg, []Any{preconditionFailure}},
	{OutOfRange, OutOfRange, msg, []Any{badRequest}},
	{Unauthenticated, Unauthenticated, msg, []Any{errInfo}},
	{PermissionDenied, PermissionDenied, msg, []Any{errInfo}},
	{Aborted, Aborted, msg, []Any{errInfo}},
	{AlreadyExists, AlreadyExists, msg, []Any{resourceInfo}},
	{ResourceExhausted, ResourceExhausted, msg, []Any{quotaFailure}},
	{Cancelled, Cancelled, msg, nil},
	{DataLoss, DataLoss, msg, []Any{debugInfo}},
	{Unknown, Unknown, msg, []Any{debugInfo}},
	{Internal, Internal, msg, []Any{debugInfo}},
	{Unimplemented, Unimplemented, msg, nil},
	{Unavailable, Unavailable, msg, []Any{debugInfo}},
	{DeadlineExceeded, DeadlineExceeded, msg, []Any{debugInfo}},
}

func TestPredefined(t *testing.T) {
	errors := []error{
		NewNotFound(msg, resourceInfo),
		NewBadRequest(msg, badRequest),
		NewFailedPrecondition(msg, preconditionFailure),
		NewOutOfRange(msg, badRequest),
		NewUnauthenticated(msg, errInfo),
		NewPermissionDenied(msg, errInfo),
		NewAborted(msg, errInfo),
		NewAlreadyExists(msg, resourceInfo),
		NewResourceExhausted(msg, quotaFailure),
		NewCancelled(msg),
		NewDataLoss(msg, debugInfo),
		NewUnknown(msg, debugInfo),
		NewInternal(msg, debugInfo),
		NewUnimplemented(msg),
		NewUnavailable(msg, debugInfo),
		NewDeadlineExceeded(msg, debugInfo),
	}

	for i, err := range errors {
		assert.Equal(t, preErrors[i], err)
	}
}
