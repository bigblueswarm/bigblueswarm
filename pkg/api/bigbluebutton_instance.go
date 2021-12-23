package api

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Create execute a create api call on the remote BigBlueButton instance
func (i *BigBlueButtonInstance) Create(params string) *CreateResponse {
	checksum := CreateChecksum(i.Secret, CreateAPIAction, params)

	body, err := i.callAPI(params, checksum)

	if err != nil {
		log.Error(fmt.Sprintf("Failed to call %s instance api create", i.URL), err)
		return nil
	}

	var response CreateResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		log.Error("Failed to unmarshal create api call body content", err)
		return nil
	}

	return &response
}

func (i *BigBlueButtonInstance) callAPI(params string, checksum *Checksum) ([]byte, error) {
	checksumValue, err := checksum.Process()
	if err != nil {
		log.Error(fmt.Sprintf("Failed to execute %s call", checksum.Action), err)
		return nil, err
	}

	url := i.URL + "/api/" + checksum.Action + "?" + params + "&checksum=" + checksumValue
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Error(fmt.Sprintf("Calling %s action on %s instance throws an exception", checksum.Action, i.URL), err)
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

// GetJoinRedirectURL compute the join redirect url
func (i *BigBlueButtonInstance) GetJoinRedirectURL(params string) (string, error) {
	checksum := CreateChecksum(i.Secret, JoinAPIAction, params)
	checksumValue, err := checksum.Process()
	if err != nil {
		log.Error("Failed to compute checksum while getting join redirect url", err)
		return "", err
	}

	return i.URL + "/api/" + JoinAPIAction + "?" + checksum.Params + "&checksum=" + checksumValue, nil
}
