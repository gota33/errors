package errors

func predefined(message string, code StatusCode, details ...Any) error {
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

func NewBadRequest(message string, detail BadRequest) error {
	return predefined(message, InvalidArgument, detail)
}

func NewFailedPrecondition(message string, detail PreconditionFailure) error {
	return predefined(message, FailedPrecondition, detail)
}

func NewOutOfRange(message string, detail BadRequest) (err error) {
	return predefined(message, OutOfRange, detail)
}

func NewUnauthenticated(message string, detail ErrorInfo) (err error) {
	return predefined(message, Unauthenticated, detail)
}

func NewPermissionDenied(message string, detail ErrorInfo) (err error) {
	return predefined(message, PermissionDenied, detail)
}

func NewAborted(message string, detail ErrorInfo) (err error) {
	return predefined(message, Aborted, detail)
}

func NewAlreadyExists(message string, detail ResourceInfo) (err error) {
	return predefined(message, AlreadyExists, detail)
}

func NewResourceExhausted(message string, detail QuotaFailure) (err error) {
	return predefined(message, ResourceExhausted, detail)
}

func NewCancelled(message string) (err error) {
	return predefined(message, Cancelled)
}

func NewDataLoss(message string, detail DebugInfo) (err error) {
	return predefined(message, DataLoss, detail)
}

func NewUnknown(message string, detail DebugInfo) (err error) {
	return predefined(message, Unknown, detail)
}

func NewInternal(message string, detail DebugInfo) (err error) {
	return predefined(message, Internal, detail)
}

func NewUnimplemented(message string) (err error) {
	return predefined(message, Unimplemented)
}

func NewUnavailable(message string, detail DebugInfo) (err error) {
	return predefined(message, Unavailable, detail)
}

func NewDeadlineExceeded(message string, detail DebugInfo) (err error) {
	return predefined(message, DeadlineExceeded, detail)
}
