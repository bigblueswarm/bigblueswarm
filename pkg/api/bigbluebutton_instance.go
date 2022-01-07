package api

import (
	"b3lb/pkg/restclient"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Create execute a create api call on the remote BigBlueButton instance
func (i *BigBlueButtonInstance) Create(params string) (*CreateResponse, error) {
	checksum := CreateChecksum(i.Secret, Create, params)

	body, err := i.callAPI(params, checksum)

	if err != nil {
		log.Error(fmt.Sprintf("Failed to call %s instance api create", i.URL), err)
		return nil, err
	}

	var response CreateResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		log.Error("Failed to unmarshal create api call body content", err)
		return nil, err
	}

	return &response, nil
}

func (i *BigBlueButtonInstance) callAPI(params string, checksum *Checksum) ([]byte, error) {
	checksumValue, err := checksum.Process()
	if err != nil {
		log.Error(fmt.Sprintf("Failed to execute %s call", checksum.Action), err)
		return nil, err
	}

	url := i.URL + "/api/" + checksum.Action + "?" + params + "&checksum=" + checksumValue
	resp, err := restclient.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Error(fmt.Sprintf("Calling %s action on %s instance throws an exception", checksum.Action, i.URL), err)
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

// GetJoinRedirectURL compute the join redirect url
func (i *BigBlueButtonInstance) GetJoinRedirectURL(params string) (string, error) {
	checksum := CreateChecksum(i.Secret, Join, params)
	checksumValue, err := checksum.Process()
	if err != nil {
		log.Error("Failed to compute checksum while getting join redirect url", err)
		return "", err
	}

	return i.URL + "/api/" + Join + "?" + checksum.Params + "&checksum=" + checksumValue, nil
}

// End execute a end api call on the remote BigBlueButton instance
func (i *BigBlueButtonInstance) End(params string) (*EndResponse, error) {
	checksum := CreateChecksum(i.Secret, End, params)

	body, err := i.callAPI(params, checksum)

	if err != nil {
		log.Error(fmt.Sprintf("Failed to call %s instance api end", i.URL), err)
		return nil, err
	}

	var response EndResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		log.Error("Failed to unmarshal end api call body content", err)
		return nil, err
	}

	return &response, nil
}

// IsMeetingRunning checks if a meeting is running on the remote Bigbluebutton instance
func (i *BigBlueButtonInstance) IsMeetingRunning(params string) (*IsMeetingsRunningResponse, error) {
	checksum := CreateChecksum(i.Secret, IsMeetingRunning, params)
	body, err := i.callAPI(params, checksum)

	if err != nil {
		log.Error(fmt.Sprintf("Failed to call %s instance api is meeting running", i.URL), err)
		return nil, err
	}

	var response IsMeetingsRunningResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		log.Error("Failed to unmarshal is meting running api call body content", err)
		return nil, err
	}

	return &response, nil
}

// GetMeetingInfo execute a get meeting info api call on the remote BigBlueButton instance
func (i *BigBlueButtonInstance) GetMeetingInfo(params string) (*GetMeetingInfoResponse, error) {
	checksum := CreateChecksum(i.Secret, GetMeetingInfo, params)
	body, err := i.callAPI(params, checksum)

	if err != nil {
		log.Error(fmt.Sprintf("Failed to call %s instance api get meeting info", i.URL), err)
		return nil, err
	}

	var response GetMeetingInfoResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		log.Error("Failed to unmarshal get meeting info api call body content", err)
		return nil, err
	}

	return &response, nil
}
