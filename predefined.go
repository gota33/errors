package errors

func predefined(cause error, code StatusCode, details ...Any) error {
	if cause == nil {
		return nil
	}

	return &annotated{
		cause:   cause,
		code:    code,
		message: cause.Error(),
		details: details,
	}
}

func WithNotFound(cause error, detail ResourceInfo) error {
	return predefined(cause, NotFound, detail)
}

func WithBadRequest(cause error, detail BadRequest) error {
	return predefined(cause, InvalidArgument, detail)
}

func WithFailedPrecondition(cause error, detail PreconditionFailure) error {
	return predefined(cause, FailedPrecondition, detail)
}

func WithOutOfRange(cause error, detail BadRequest) (err error) {
	return predefined(cause, OutOfRange, detail)
}

func WithUnauthenticated(cause error, detail ErrorInfo) (err error) {
	return predefined(cause, Unauthenticated, detail)
}

func WithPermissionDenied(cause error, detail ErrorInfo) (err error) {
	return predefined(cause, PermissionDenied, detail)
}

func WithAborted(cause error, detail ErrorInfo) (err error) {
	return predefined(cause, Aborted, detail)
}

func WithAlreadyExists(cause error, detail ResourceInfo) (err error) {
	return predefined(cause, AlreadyExists, detail)
}

func WithResourceExhausted(cause error, detail QuotaFailure) (err error) {
	return predefined(cause, ResourceExhausted, detail)
}

func WithCancelled(cause error) (err error) {
	return predefined(cause, Cancelled)
}

func WithDataLoss(cause error, detail DebugInfo) (err error) {
	return predefined(cause, DataLoss, detail)
}

func WithUnknown(cause error, detail DebugInfo) (err error) {
	return predefined(cause, Unknown, detail)
}

func WithInternal(cause error, detail DebugInfo) (err error) {
	return predefined(cause, Internal, detail)
}

func WithUnimplemented(cause error) (err error) {
	return predefined(cause, Unimplemented)
}

func WithUnavailable(cause error, detail DebugInfo) (err error) {
	return predefined(cause, Unavailable, detail)
}

func WithDeadlineExceeded(cause error, detail DebugInfo) (err error) {
	return predefined(cause, DeadlineExceeded, detail)
}
