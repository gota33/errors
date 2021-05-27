package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const msg = "msg"

var rootErr = fmt.Errorf(msg)

var preErrors = []*annotated{
	{rootErr, NotFound, msg, []Any{resourceInfo}},
	{rootErr, InvalidArgument, msg, []Any{badRequest}},
	{rootErr, FailedPrecondition, msg, []Any{preconditionFailure}},
	{rootErr, OutOfRange, msg, []Any{badRequest}},
	{rootErr, Unauthenticated, msg, []Any{errInfo}},
	{rootErr, PermissionDenied, msg, []Any{errInfo}},
	{rootErr, Aborted, msg, []Any{errInfo}},
	{rootErr, AlreadyExists, msg, []Any{resourceInfo}},
	{rootErr, ResourceExhausted, msg, []Any{quotaFailure}},
	{rootErr, Cancelled, msg, nil},
	{rootErr, DataLoss, msg, []Any{debugInfo}},
	{rootErr, Unknown, msg, []Any{debugInfo}},
	{rootErr, Internal, msg, []Any{debugInfo}},
	{rootErr, Unimplemented, msg, nil},
	{rootErr, Unavailable, msg, []Any{debugInfo}},
	{rootErr, DeadlineExceeded, msg, []Any{debugInfo}},
}

func TestPredefined(t *testing.T) {
	errors := []error{
		WithNotFound(rootErr, resourceInfo),
		WithBadRequest(rootErr, badRequest),
		WithFailedPrecondition(rootErr, preconditionFailure),
		WithOutOfRange(rootErr, badRequest),
		WithUnauthenticated(rootErr, errInfo),
		WithPermissionDenied(rootErr, errInfo),
		WithAborted(rootErr, errInfo),
		WithAlreadyExists(rootErr, resourceInfo),
		WithResourceExhausted(rootErr, quotaFailure),
		WithCancelled(rootErr),
		WithDataLoss(rootErr, debugInfo),
		WithUnknown(rootErr, debugInfo),
		WithInternal(rootErr, debugInfo),
		WithUnimplemented(rootErr),
		WithUnavailable(rootErr, debugInfo),
		WithDeadlineExceeded(rootErr, debugInfo),
	}

	for i, err := range errors {
		assert.Equal(t, preErrors[i], err)
	}
}

func TestNil(t *testing.T) {
	err := WithUnknown(nil, DebugInfo{})
	assert.NoError(t, err)
}
