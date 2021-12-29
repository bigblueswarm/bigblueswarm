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

// returnCodes represents the api return code
type returnCodes struct {
	Success string
	Failed  string
}

// ReturnCodes returns a struct containing the api return codes
func ReturnCodes() *returnCodes {
	return &returnCodes{
		Success: "SUCCESS",
		Failed:  "FAILED",
	}
}

// messageKeys represents the api message key
type messageKeys struct {
	ValidationError    string
	DuplicationWarning string
	NotFound           string
}

// MessageKeys return a struct containing the api message keys
func MessageKeys() *messageKeys {
	return &messageKeys{
		ValidationError:    "validationError",
		DuplicationWarning: "duplicationWarning",
		NotFound:           "notFound",
	}
}

// messages represents the api messages
type messages struct {
	EmptyMeetingID     string
	EmptyMeetingName   string
	DuplicationWarning string
	NotFound           string
}

// Messages returns a struct containing the api messages
func Messages() *messages {
	return &messages{
		EmptyMeetingID:     "You must provide a meeting ID",
		EmptyMeetingName:   "You must provide a meeting name",
		DuplicationWarning: "This conference was already in existence and may currently be in progress.",
		NotFound:           "A meeting with that ID does not exist",
	}
}
