package errors

func predefined(message string, code StatusCode, details ...Detail) error {
	return &annotated{
		cause:   code,
		code:    code,
		message: message,
		details: details,
	}
}

func NewNotFound(message string, detail ResourceInfo) error {
	return predefined(message, NotFound, detail)
}

func NewBadRequest(message string, violations ...FieldViolation) error {
	return predefined(message, InvalidArgument, BadRequest{violations})
}

func NewFailedPrecondition(message string, violations ...TypedViolation) error {
	return predefined(message, FailedPrecondition, PreconditionFailure{violations})
}

func NewOutOfRange(message string, violations ...FieldViolation) (err error) {
	return predefined(message, OutOfRange, BadRequest{violations})
}

func NewUnauthenticated(message string, info ErrorInfo) (err error) {
	return predefined(message, Unauthenticated, info)
}

func NewPermissionDenied(message string, info ErrorInfo) (err error) {
	return predefined(message, PermissionDenied, info)
}

func NewAborted(message string, info ErrorInfo) (err error) {
	return predefined(message, Aborted, info)
}

func NewAlreadyExists(message string, info ResourceInfo) (err error) {
	return predefined(message, AlreadyExists, info)
}

func NewResourceExhausted(message string, violations ...Violation) (err error) {
	return predefined(message, ResourceExhausted, QuotaFailure{violations})
}

func NewCancelled(message string) (err error) {
	return predefined(message, Cancelled)
}

func NewDataLoss(message string, info DebugInfo) (err error) {
	return predefined(message, DataLoss, info)
}

func NewUnknown(message string, info DebugInfo) (err error) {
	return predefined(message, Unknown, info)
}

func NewInternal(message string, info DebugInfo) (err error) {
	return predefined(message, Internal, info)
}

func NewUnimplemented(message string) (err error) {
	return predefined(message, Unimplemented)
}

func NewUnavailable(message string, info DebugInfo) (err error) {
	return predefined(message, Unavailable, info)
}

func NewDeadlineExceeded(message string, info DebugInfo) (err error) {
	return predefined(message, DeadlineExceeded, info)
}
