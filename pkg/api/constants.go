package api

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
	ValidationError       string
	DuplicationWarning    string
	NotFound              string
	SendEndMeetingRequest string
}

// MessageKeys return a struct containing the api message keys
func MessageKeys() *Keys {
	return &Keys{
		ValidationError:       "validationError",
		DuplicationWarning:    "duplicationWarning",
		NotFound:              "notFound",
		SendEndMeetingRequest: "sentEndMeetingRequest",
	}
}

// MessageValues represents the api messages
type MessageValues struct {
	EmptyMeetingID     string
	EmptyMeetingName   string
	DuplicationWarning string
	NotFound           string
	EndMeeting         string
	InvalidModeratorPW string
}

// Messages returns a struct containing the api messages
func Messages() *MessageValues {
	return &MessageValues{
		EmptyMeetingID:     "You must provide a meeting ID",
		EmptyMeetingName:   "You must provide a meeting name",
		DuplicationWarning: "This conference was already in existence and may currently be in progress.",
		NotFound:           "A meeting with that ID does not exist",
		EndMeeting:         "A request to end the meeting was sent. Please wait a few seconds, and then use the getMeetingInfo or isMeetingRunning API calls to verify that it was ended.",
		InvalidModeratorPW: "Provided moderator password is incorrect",
	}
}
