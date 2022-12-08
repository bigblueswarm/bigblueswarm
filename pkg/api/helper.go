// Package api manage the bigbluebutton api and communication between bigblueswarm and bigbluebutton instances
package api

// CreateHealthCheck returns a health check response
func CreateHealthCheck() *HealthCheck {
	return &HealthCheck{
		ReturnCode: "SUCCESS",
		Version:    Version,
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

// CreateJSONError returns an error response given a message key and message
func CreateJSONError(key string, message string) *JSONResponse {
	return &JSONResponse{
		Response{
			ReturnCode: "FAILED",
			MessageKey: key,
			Message:    message,
		},
	}
}

// CreateChecksum returns a checksum given a secret, action and params
func CreateChecksum(secret string, action string, params string) *Checksum {
	return &Checksum{
		Secret: secret,
		Action: action,
		Params: params,
	}
}
