package api

import "encoding/xml"

type HealtCheck struct {
	XMLName     xml.Name `xml:"response"`
	Return_code string   `xml:"returncode"`
	Version     string   `xml:"version"`
}

type ChecksumError struct {
	XMLName    xml.Name `xml:"response"`
	ReturnCode string   `xml:"returncode"`
	MessageKey string   `xml:"messageKey"`
	Message    string   `xml:"message"`
}

func DefaultChecksumError() *ChecksumError {
	return &ChecksumError{
		ReturnCode: "FAILED",
		MessageKey: "checksumError",
		Message:    "You did not pass the checksum security check",
	}
}
