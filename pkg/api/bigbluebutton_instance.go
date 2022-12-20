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

func (i *BigBlueButtonInstance) getLogger(action string, params string) *log.Entry {
	return log.WithFields(log.Fields{
		"instance": i.URL,
		"action":   action,
		"params":   params,
	})
}

// Create execute a create api call on the remote BigBlueButton instance
func (i *BigBlueButtonInstance) Create(params string) (*CreateResponse, error) {
	logger := i.getLogger(Create, params)
	response, err := i.api(Create, params)

	if err != nil {
		logger.Error("api call to create method throws an error")
		return nil, err
	}

	if creation, ok := response.(*CreateResponse); ok {
		return creation, nil
	}

	e := errors.New("failed to cast api response to CreateResponse")
	log.Error(e)
	return nil, e
}

func (i *BigBlueButtonInstance) callAPI(checksum *Checksum) ([]byte, error) {
	logger := i.getLogger(checksum.Action, checksum.Params)
	checksumValue, err := checksum.Process()
	if err != nil {
		logger.Error(fmt.Sprintf("failed to execute %s call. The checksum process failed", checksum.Action), err)
		return nil, err
	}

	url := i.URL + "/api/" + checksum.Action + "?" + checksum.Params + "&checksum=" + checksumValue
	resp, err := restclient.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("calling %s action on %s instance throws an exception", checksum.Action, i.URL), err)
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
	logger := i.getLogger(action, params)
	checksum := CreateChecksum(i.Secret, action, params)

	body, err := i.callAPI(checksum)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to call %s instance %s api", i.URL, action), err)
		return nil, err
	}

	dataType := actionMapper(action)
	result := reflect.New(dataType).Interface()
	if err := unMarshall(action, body, &result); err != nil {
		logger.Error(fmt.Sprintf("failed to unmarshal %s api call body content", action), err)
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
	logger := i.getLogger(Join, params)
	checksum := CreateChecksum(i.Secret, Join, params)
	checksumValue, err := checksum.Process()
	if err != nil {
		logger.Error("failed to compute checksum while getting join redirect url", err)
		return "", err
	}

	return i.URL + "/api/" + Join + "?" + checksum.Params + "&checksum=" + checksumValue, nil
}

// Join execute a join api call on the remote BigBlueButton instance
func (i *BigBlueButtonInstance) Join(params string) (*JoinRedirectResponse, error) {
	logger := i.getLogger(Join, params)
	response, err := i.api(Join, params)

	if err != nil {
		logger.Error("api call to Join api failed", err)
		return nil, err
	}

	if join, ok := response.(*JoinRedirectResponse); ok {
		return join, nil
	}

	e := errors.New("failed to cast api response to JoinRedirectResponse")
	logger.Error(e)
	return nil, e
}

// End execute a end api call on the remote BigBlueButton instance
func (i *BigBlueButtonInstance) End(params string) (*EndResponse, error) {
	logger := i.getLogger(End, params)
	response, err := i.api(End, params)

	if err != nil {
		logger.Error("api call to End api failed", err)
		return nil, err
	}

	if end, ok := response.(*EndResponse); ok {
		return end, nil
	}

	e := errors.New("failed to cast api response to EndResponse")
	logger.Error(e)
	return nil, e
}

// IsMeetingRunning checks if a meeting is running on the remote Bigbluebutton instance
func (i *BigBlueButtonInstance) IsMeetingRunning(params string) (*IsMeetingsRunningResponse, error) {
	logger := i.getLogger(IsMeetingRunning, params)
	response, err := i.api(IsMeetingRunning, params)

	if err != nil {
		logger.Error("api call to IsMeetingRunning api failed", err)
		return nil, err
	}

	if running, ok := response.(*IsMeetingsRunningResponse); ok {
		return running, nil
	}

	e := errors.New("failed to cast api response to IsMeetingsRunningResponse")
	logger.Error(e)
	return nil, e
}

// GetMeetingInfo execute a get meeting info api call on the remote BigBlueButton instance
func (i *BigBlueButtonInstance) GetMeetingInfo(params string) (*GetMeetingInfoResponse, error) {
	logger := i.getLogger(GetMeetingInfo, params)
	response, err := i.api(GetMeetingInfo, params)

	if err != nil {
		logger.Error("api call to GetMeetingInfo api failed", err)
		return nil, err
	}

	if meeting, ok := response.(*GetMeetingInfoResponse); ok {
		return meeting, nil
	}

	e := errors.New("failed to cast api response to GetMeetingInfoResponse")
	logger.Error(e)
	return nil, e
}

// GetMeetings execute a get meetings api call on the remote BigBlueButton instance
func (i *BigBlueButtonInstance) GetMeetings() (*GetMeetingsResponse, error) {
	logger := i.getLogger(GetMeetings, "")
	response, err := i.api(GetMeetings, "")

	if err != nil {
		logger.Error("api call to GetMeetings api failed", err)
		return nil, err
	}

	if meetings, ok := response.(*GetMeetingsResponse); ok {
		return meetings, nil
	}

	e := errors.New("failed to cast api response to GetMeetingsResponse")
	logger.Error(e)
	return nil, e
}

// GetRecordings perform a get recordings api call on the remote BigBlueButton instance
func (i *BigBlueButtonInstance) GetRecordings(params string) (*GetRecordingsResponse, error) {
	logger := i.getLogger(GetRecordings, params)
	response, err := i.api(GetRecordings, params)

	if err != nil {
		logger.Error("api call to GetRecordings api failed", err)
		return nil, err
	}

	if recordings, ok := response.(*GetRecordingsResponse); ok {
		return recordings, nil
	}

	e := errors.New("failed to cast api response to GetRecordingsResponse")
	logger.Error(e)
	return nil, e
}

// UpdateRecordings perform a update recordings api call on the remote BigBlueButton instance
func (i *BigBlueButtonInstance) UpdateRecordings(params string) (*UpdateRecordingsResponse, error) {
	logger := i.getLogger(UpdateRecordings, params)
	response, err := i.api(UpdateRecordings, params)

	if err != nil {
		logger.Error("api call to UpdateRecordings api failed", err)
		return nil, err
	}

	if recordings, ok := response.(*UpdateRecordingsResponse); ok {
		return recordings, nil
	}

	e := errors.New("failed to cast api response to UpdateRecordingsResponse")
	logger.Error(e)
	return nil, e
}

// DeleteRecordings perform a delete recordings api call on the remote BigBlueButton instance
func (i *BigBlueButtonInstance) DeleteRecordings(params string) (*DeleteRecordingsResponse, error) {
	logger := i.getLogger(DeleteRecordings, params)
	response, err := i.api(DeleteRecordings, params)

	if err != nil {
		logger.Error("api call to DeleteRecordings api failed", err)
		return nil, err
	}

	if deletion, ok := response.(*DeleteRecordingsResponse); ok {
		return deletion, nil
	}

	e := errors.New("failed to cast api response to DeleteRecordingsResponse")
	logger.Error(e)
	return nil, e
}

// PublishRecordings perform a publish recordings api call on the remote BigBlueButton instance
func (i *BigBlueButtonInstance) PublishRecordings(params string) (*PublishRecordingsResponse, error) {
	logger := i.getLogger(PublishRecordings, params)
	response, err := i.api(PublishRecordings, params)

	if err != nil {
		logger.Error("api call to PublishRecordings api failed", err)
		return nil, err
	}

	if publish, ok := response.(*PublishRecordingsResponse); ok {
		return publish, nil
	}

	e := errors.New("failed to cast api response to PublishRecordingsResponse")
	logger.Error(e)
	return nil, e
}

// GetRecordingTextTracks perform a get recording text tracks api call on the remote BigBlueButton instance
func (i *BigBlueButtonInstance) GetRecordingTextTracks(params string) (*GetRecordingsTextTracksResponse, error) {
	logger := i.getLogger(GetRecordingsTextTracks, params)
	response, err := i.api(GetRecordingsTextTracks, params)

	if err != nil {
		logger.Error("api call to GetRecordingsTextTracks api failed", err)
		return nil, err
	}

	if tracks, ok := response.(*GetRecordingsTextTracksResponse); ok {
		return tracks, nil
	}

	e := errors.New("failed to cast api response to GetRecordingsTextTracksResponse")
	logger.Error(e)
	return nil, e
}

// Redirect redirect provided context to instance action
func (i *BigBlueButtonInstance) Redirect(c *gin.Context, action string, parameters string) {
	logger := i.getLogger(action, parameters)
	checksum := CreateChecksum(i.Secret, action, parameters)
	checksumValue, err := checksum.Process()
	if err != nil {
		logger.Error("failed to redirect", err)
		c.XML(http.StatusOK, CreateError(MessageKeys().NotFound, Messages().NotFound))
		return
	}

	url := i.URL + "/api/" + action + "?" + checksum.Params + "&checksum=" + checksumValue
	c.Writer.Header().Set("Location", url)
	c.AbortWithStatus(http.StatusFound)
}
