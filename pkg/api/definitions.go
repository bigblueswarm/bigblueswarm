package api

import "encoding/xml"

type HealtCheck struct {
	XMLName     xml.Name `xml:"response"`
	Return_code string   `xml:"returncode"`
	Version     string   `xml:"version"`
}

type ChecksumError struct {
	XMLName     xml.Name `xml:"response"`
	Return_code string   `xml:"returncode"`
	Message_key string   `xml:"messageKey"`
	Message     string   `xml:"message"`
}
