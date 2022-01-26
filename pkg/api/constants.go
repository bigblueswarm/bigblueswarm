package api

// Version is the version of the api
const Version = "2.0"

// BigBlueButton is the main endpoint for the API
const BigBlueButton = "bigbluebutton"

// API is the sub-endpoint for the API
const API = "api"

// Create is the sub-endpoint for creating a meeting
const Create = "create"

// GetMeetings is the sub-endpoint for getting a list of meetings
const GetMeetings = "getMeetings"

// Join is the sub-endpoint for joining a meeting
const Join = "join"

// End is the sub-endpoint for ending a meeting
const End = "end"

// IsMeetingRunning is the sub-endpoint for checking if a meeting is running
const IsMeetingRunning = "isMeetingRunning"

// GetMeetingInfo is the sub-endpoint for getting a meeting info
const GetMeetingInfo = "getMeetingInfo"

// GetRecordings is the sub-endpoint for getting a list of recordings
const GetRecordings = "getRecordings"

// UpdateRecordings is the sub-endpoint for updating a recording
const UpdateRecordings = "updateRecordings"

// DeleteRecordings is the sub-endpoint for deleting a recording
const DeleteRecordings = "deleteRecordings"

// PublishRecordings is the sub-endpoint for publishing a recording
const PublishRecordings = "publishRecordings"

// Codes represents the api return code
type Codes struct {
	Success string
	Failed  string
}

// ReturnCodes returns a struct containing the api return codes
func ReturnCodes() *Codes {
	return &Codes{
		Success: "SUCCESS",
		Failed:  "FAILED",
	}
}

// Keys represents the api message key
type Keys struct {
	ValidationError          string
	DuplicationWarning       string
	NotFound                 string
	SendEndMeetingRequest    string
	NoRecordings             string
	MissingRecordIDParameter string
}

// MessageKeys return a struct containing the api message keys
func MessageKeys() *Keys {
	return &Keys{
		ValidationError:          "validationError",
		DuplicationWarning:       "duplicationWarning",
		NotFound:                 "notFound",
		SendEndMeetingRequest:    "sentEndMeetingRequest",
		NoRecordings:             "noRecordings",
		MissingRecordIDParameter: "missingParamRecordID",
	}
}

// MessageValues represents the api messages
type MessageValues struct {
	EmptyMeetingID           string
	EmptyMeetingName         string
	DuplicationWarning       string
	NotFound                 string
	EndMeeting               string
	InvalidModeratorPW       string
	NoRecordings             string
	MissingRecordIDParameter string
	RecordingNotFound        string
}

// Messages returns a struct containing the api messages
func Messages() *MessageValues {
	return &MessageValues{
		EmptyMeetingID:           "You must provide a meeting ID",
		EmptyMeetingName:         "You must provide a meeting name",
		DuplicationWarning:       "This conference was already in existence and may currently be in progress.",
		NotFound:                 "A meeting with that ID does not exist",
		EndMeeting:               "A request to end the meeting was sent. Please wait a few seconds, and then use the getMeetingInfo or isMeetingRunning API calls to verify that it was ended.",
		InvalidModeratorPW:       "Provided moderator password is incorrect",
		NoRecordings:             "There are no recordings for the meeting(s).",
		MissingRecordIDParameter: "You must specify a recordID.",
		RecordingNotFound:        "We could not find recordings",
	}
}
