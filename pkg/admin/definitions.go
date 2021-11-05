package admin

// BigBlueButtonInstance represents a REST admin Bigbluebutton instance. It contains the  server URL and the server secret.
type BigBlueButtonInstance struct {
	URL    string `json:"url"`
	Secret string `json:"secret"`
}
