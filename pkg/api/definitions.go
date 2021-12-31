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

// Attendee represents a Bigbluebutton attendee
type Attendee struct {
	UserID          string `xml:"userID"`
	FullName        string `xml:"fullName"`
	Role            string `xml:"role"`
	IsPresenter     bool   `xml:"isPresenter"`
	IsListeningOnly bool   `xml:"isListeningOnly"`
	HasJoinedVoice  bool   `xml:"hasJoinedVoice"`
	HasVideo        bool   `xml:"hasVideo"`
	ClientType      string `xml:"clientType"`
}

// GetMeetingInfoResponse represents the Bigbluebutton getMeetingInfo API response type
type GetMeetingInfoResponse struct {
	XMLName               xml.Name   `xml:"response"`
	ReturnCode            string     `xml:"returncode"`
	InternalMeetingID     string     `xml:"internalMeetingID"`
	MeetingID             string     `xml:"meetingID"`
	CreateTime            string     `xml:"createTime"`
	CreateDate            string     `xml:"createDate"`
	VoiceBridge           string     `xml:"voiceBridge"`
	DialNumber            string     `xml:"dialNumber"`
	AttendeePW            string     `xml:"attendeePW"`
	ModeratorPW           string     `xml:"moderatorPW"`
	Running               bool       `xml:"running"`
	Duration              int        `xml:"duration"`
	HasUserJoined         string     `xml:"hasUserJoined"`
	Recording             bool       `xml:"recording"`
	HasBeenForciblyEnded  bool       `xml:"hasBeenForciblyEnded"`
	StartTime             int        `xml:"startTime"`
	EndTime               int        `xml:"endTime"`
	ParticipantCount      int        `xml:"participantCount"`
	ListenerCount         int        `xml:"listenerCount"`
	VoiceParticipantCount int        `xml:"voiceParticipantCount"`
	VideoCount            int        `xml:"videoCount"`
	MaxUsers              int        `xml:"maxUsers"`
	ModeratorCount        int        `xml:"moderatorCount"`
	Attendees             []Attendee `xml:"attendees>attendee"`
	MetaData              struct {
		Inner []byte `xml:",innerxml"`
	} `xml:"metadata"`
	IsBreakout bool `xml:"isBreakout"`
}
