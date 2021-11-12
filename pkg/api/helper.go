package api

// CreateHealthCheck returns a health check response
func CreateHealthCheck() *HealthCheck {
	return &HealthCheck{
		ReturnCode: "SUCCESS",
		Version:    "2.0",
	}
}

// DefaultChecksumError returns a default checksum error
func DefaultChecksumError() *Error {
	return CreateError("checksumError", "You did not pass the checksum security check")
}

// CreateError returns an error response given a message key and message
func CreateError(key string, message string) *Error {
	return &Error{
		Response{
			ReturnCode: "FAILED",
			MessageKey: key,
			Message:    message,
		},
	}
}
