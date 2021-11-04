package api

import "encoding/xml"

// Represents the healthcheck response
type HealtCheck struct {
	XMLName    xml.Name `xml:"response"`
	ReturnCode string   `xml:"returncode"`
	Version    string   `xml:"version"`
}

// Represents the checksum error response
type ChecksumError struct {
	XMLName    xml.Name `xml:"response"`
	ReturnCode string   `xml:"returncode"`
	MessageKey string   `xml:"messageKey"`
	Message    string   `xml:"message"`
}

// Returns a default checksum error
func DefaultChecksumError() *ChecksumError {
	return &ChecksumError{
		ReturnCode: "FAILED",
		MessageKey: "checksumError",
		Message:    "You did not pass the checksum security check",
	}
}
