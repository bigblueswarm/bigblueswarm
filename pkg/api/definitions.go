package api

import (
	"encoding/xml"
)

// Checksum in BigBlueButton authentication system represents an action name, all parameters and a secret concatenated in a single string that is hashed by SHA1.
type Checksum struct {
	Secret string
	Action string
	Params string
}

// BigBlueButtonInstance represents a REST admin Bigbluebutton instance. It contains the  server URL and the server secret.
type BigBlueButtonInstance struct {
	URL    string `json:"url"`
	Secret string `json:"secret"`
}

// HealthCheck represents the healthcheck response
type HealthCheck struct {
	XMLName    xml.Name `xml:"response"`
	ReturnCode string   `xml:"returncode"`
	Version    string   `xml:"version"`
}

// Error represents the error response
type Error struct {
	Response
}

// Response represents the basic api response
type Response struct {
	XMLName    xml.Name `xml:"response"`
	ReturnCode string   `xml:"returncode"`
	MessageKey string   `xml:"messageKey"`
	Message    string   `xml:"message"`
}

// CreateResponse represents the Bigbluebutton create API response type
type CreateResponse struct {
	Response
	MeetingID            string `xml:"meetingID"`
	InternalMeetingID    string `xml:"internalMetingID"`
	ParentMeetingID      string `xml:"parentMetingID"`
	AttendeePW           string `xml:"attendePW"`
	ModeratorPW          string `xml:"moderatorPW"`
	CreateTime           string `xml:"createTime"`
	VoiceBridge          string `xml:"voiceBridge"`
	DialNumber           string `xml:"dialNumber"`
	CreateDate           string `xml:"createDate"`
	HasUserJoined        string `xml:"hasUserJoined"`
	Duration             string `xml:"duration"`
	HasBeenForciblyEnded string `xml:"hasBeenForciblyEnded"`
}

// EndResponse represents the Bigbluebutton end API response type
type EndResponse struct {
	Response
}

// IsMeetingsRunningResponse represents the Bigbluebutton isMeetingRunning API response type
type IsMeetingsRunningResponse struct {
	XMLName    xml.Name `xml:"response"`
	ReturnCode string   `xml:"returncode"`
	Running    bool     `xml:"running"`
}
