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

// MeetingInfo represents the Bigbluebutton meeting info API object
type MeetingInfo struct {
	MeetingName           string     `xml:"meetingName"`
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

// GetMeetingInfoResponse represents the Bigbluebutton getMeetingInfo API response type
type GetMeetingInfoResponse struct {
	XMLName    xml.Name `xml:"response"`
	ReturnCode string   `xml:"returncode"`
	MeetingInfo
}

// GetMeetingsResponse represents the Bigbluebutton getMeetings API response type
type GetMeetingsResponse struct {
	XMLName    xml.Name      `xml:"response"`
	ReturnCode string        `xml:"returncode"`
	Meetings   []MeetingInfo `xml:"meetings>meeting"`
}

// JoinRedirectResponse represents the BigBlueButton join API response type when query parameter `redirect=false` is set
type JoinRedirectResponse struct {
	Response
	MeetingID    string `xml:"meeting_id"`
	UserID       string `xml:"user_id"`
	AuthToken    string `xml:"auth_token"`
	SessionToken string `xml:"session_token"`
	URL          string `xml:"url"`
}

// Recording represents the BigBlueButton recording API object
type Recording struct {
	XMLName           xml.Name `xml:"recording"`
	RecordID          string   `xml:"recordID"`
	MeetingID         string   `xml:"meetingID"`
	InternalMeetingID string   `xml:"internalMeetingID"`
	Name              string   `xml:"name"`
	IsBreakout        bool     `xml:"isBreakout"`
	Published         bool     `xml:"published"`
	State             string   `xml:"state"`
	StartTime         int      `xml:"startTime"`
	EndTime           int      `xml:"endTime"`
	Participants      int      `xml:"participants"`
	MetaData          struct {
		Inner []byte `xml:",innerxml"`
	} `xml:"metadata"`
	Playback struct {
		Inner []byte `xml:",innerxml"`
	}
}

// GetRecordingsResponse represents the Bigbluebutton getRecordings API response type
type GetRecordingsResponse struct {
	Response
	Recordings []Recording `xml:"recordings>recording"`
}

// UpdateRecordingsResponse represents the Bigbluebutton updateRecordings API response type
type UpdateRecordingsResponse struct {
	XMLName    xml.Name `xml:"response"`
	ReturnCode string   `xml:"returncode"`
	Updated    bool     `xml:"updated"`
}

// DeleteRecordingsResponse represents the Bigbluebutton deleteRecordings API response type
type DeleteRecordingsResponse struct {
	XMLName    xml.Name `xml:"response"`
	ReturnCode string   `xml:"returncode"`
	Deleted    bool     `xml:"deleted"`
}
