// Package api manage the bigbluebutton api and communication between bigblueswarm and bigbluebutton instances
package api

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/bigblueswarm/bigblueswarm/v2/pkg/restclient"
	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

// Create execute a create api call on the remote BigBlueButton instance
func (i *BigBlueButtonInstance) Create(params string) (*CreateResponse, error) {
	response, err := i.api(Create, params)

	if err != nil {
		return nil, err
	}

	if creation, ok := response.(*CreateResponse); ok {
		return creation, nil
	}

	return nil, errors.New("Failed to cast api response to CreateResponse")
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

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func unMarshall(action string, body []byte, result *interface{}) error {
	switch action {
	case GetRecordingsTextTracks:
		return json.Unmarshal(body, result)
	default:
		return xml.Unmarshal(body, result)
	}
}

func (i *BigBlueButtonInstance) api(action string, params string) (interface{}, error) {
	checksum := CreateChecksum(i.Secret, action, params)

	body, err := i.callAPI(params, checksum)

	if err != nil {
		log.Error(fmt.Sprintf("Failed to call %s instance %s api", i.URL, action), err)
		return nil, err
	}

	dataType := actionMapper(action)
	result := reflect.New(dataType).Interface()
	if err := unMarshall(action, body, &result); err != nil {
		log.Error(fmt.Sprintf("Failed to unmarshal %s api call body content", action), err)
		return nil, err
	}

	return result, nil
}

func actionMapper(action string) reflect.Type {
	switch action {
	case Create:
		return reflect.TypeOf(CreateResponse{})
	case Join:
		return reflect.TypeOf(JoinRedirectResponse{})
	case End:
		return reflect.TypeOf(EndResponse{})
	case IsMeetingRunning:
		return reflect.TypeOf(IsMeetingsRunningResponse{})
	case GetMeetingInfo:
		return reflect.TypeOf(GetMeetingInfoResponse{})
	case GetMeetings:
		return reflect.TypeOf(GetMeetingsResponse{})
	case GetRecordings:
		return reflect.TypeOf(GetRecordingsResponse{})
	case UpdateRecordings:
		return reflect.TypeOf(UpdateRecordingsResponse{})
	case DeleteRecordings:
		return reflect.TypeOf(DeleteRecordingsResponse{})
	case PublishRecordings:
		return reflect.TypeOf(PublishRecordingsResponse{})
	case GetRecordingsTextTracks:
		return reflect.TypeOf(GetRecordingsTextTracksResponse{})
	default:
		return nil
	}
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

// Join execute a join api call on the remote BigBlueButton instance
func (i *BigBlueButtonInstance) Join(params string) (*JoinRedirectResponse, error) {
	response, err := i.api(Join, params)

	if err != nil {
		return nil, err
	}

	if join, ok := response.(*JoinRedirectResponse); ok {
		return join, nil
	}

	return nil, errors.New("failed to cast api response to JoinRedirectResponse")
}

// End execute a end api call on the remote BigBlueButton instance
func (i *BigBlueButtonInstance) End(params string) (*EndResponse, error) {
	response, err := i.api(End, params)

	if err != nil {
		return nil, err
	}

	if end, ok := response.(*EndResponse); ok {
		return end, nil
	}

	return nil, errors.New("failed to cast api response to EndResponse")
}

// IsMeetingRunning checks if a meeting is running on the remote Bigbluebutton instance
func (i *BigBlueButtonInstance) IsMeetingRunning(params string) (*IsMeetingsRunningResponse, error) {
	response, err := i.api(IsMeetingRunning, params)

	if err != nil {
		return nil, err
	}

	if running, ok := response.(*IsMeetingsRunningResponse); ok {
		return running, nil
	}

	return nil, errors.New("failed to cast api response to IsMeetingsRunningResponse")
}

// GetMeetingInfo execute a get meeting info api call on the remote BigBlueButton instance
func (i *BigBlueButtonInstance) GetMeetingInfo(params string) (*GetMeetingInfoResponse, error) {
	response, err := i.api(GetMeetingInfo, params)

	if err != nil {
		return nil, err
	}

	if meeting, ok := response.(*GetMeetingInfoResponse); ok {
		return meeting, nil
	}

	return nil, errors.New("failed to cast api response to GetMeetingInfoResponse")
}

// GetMeetings execute a get meetings api call on the remote BigBlueButton instance
func (i *BigBlueButtonInstance) GetMeetings() (*GetMeetingsResponse, error) {
	response, err := i.api(GetMeetings, "")

	if err != nil {
		return nil, err
	}

	if meetings, ok := response.(*GetMeetingsResponse); ok {
		return meetings, nil
	}

	return nil, errors.New("failed to cast api response to GetMeetingsResponse")
}

// GetRecordings perform a get recordings api call on the remote BigBlueButton instance
func (i *BigBlueButtonInstance) GetRecordings(params string) (*GetRecordingsResponse, error) {
	response, err := i.api(GetRecordings, params)

	if err != nil {
		return nil, err
	}

	if recordings, ok := response.(*GetRecordingsResponse); ok {
		return recordings, nil
	}

	return nil, errors.New("failed to cast api response to GetRecordingsResponse")
}

// UpdateRecordings perform a update recordings api call on the remote BigBlueButton instance
func (i *BigBlueButtonInstance) UpdateRecordings(params string) (*UpdateRecordingsResponse, error) {
	response, err := i.api(UpdateRecordings, params)

	if err != nil {
		return nil, err
	}

	if recordings, ok := response.(*UpdateRecordingsResponse); ok {
		return recordings, nil
	}

	return nil, errors.New("failed to cast api response to UpdateRecordingsResponse")
}

// DeleteRecordings perform a delete recordings api call on the remote BigBlueButton instance
func (i *BigBlueButtonInstance) DeleteRecordings(params string) (*DeleteRecordingsResponse, error) {
	response, err := i.api(DeleteRecordings, params)

	if err != nil {
		return nil, err
	}

	if deletion, ok := response.(*DeleteRecordingsResponse); ok {
		return deletion, nil
	}

	return nil, errors.New("failed to cast api response to DeleteRecordingsResponse")
}

// PublishRecordings perform a publish recordings api call on the remote BigBlueButton instance
func (i *BigBlueButtonInstance) PublishRecordings(params string) (*PublishRecordingsResponse, error) {
	response, err := i.api(PublishRecordings, params)

	if err != nil {
		return nil, err
	}

	if publish, ok := response.(*PublishRecordingsResponse); ok {
		return publish, nil
	}

	return nil, errors.New("failed to cast api response to PublishRecordingsResponse")
}

// GetRecordingTextTracks perform a get recording text tracks api call on the remote BigBlueButton instance
func (i *BigBlueButtonInstance) GetRecordingTextTracks(params string) (*GetRecordingsTextTracksResponse, error) {
	response, err := i.api(GetRecordingsTextTracks, params)

	if err != nil {
		return nil, err
	}

	if tracks, ok := response.(*GetRecordingsTextTracksResponse); ok {
		return tracks, nil
	}

	return nil, errors.New("failed to cast api response to GetRecordingsTextTracksResponse")
}

// Redirect redirect provided context to instance action
func (i *BigBlueButtonInstance) Redirect(c *gin.Context, action string, parameters string) {
	checksum := CreateChecksum(i.Secret, action, parameters)
	checksumValue, err := checksum.Process()
	if err != nil {
		c.XML(http.StatusOK, CreateError(MessageKeys().NotFound, Messages().NotFound))
		return
	}

	url := i.URL + "/api/" + action + "?" + checksum.Params + "&checksum=" + checksumValue
	c.Writer.Header().Set("Location", url)
	c.AbortWithStatus(http.StatusFound)
}
