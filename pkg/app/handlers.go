package app

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/SLedunois/b3lb/pkg/api"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// HealthCheck handler returns an health check response
func (s *Server) HealthCheck(c *gin.Context) {
	c.XML(http.StatusOK, api.CreateHealthCheck())
}

// GetMeetings handler returns the getMeetings API. See https://docs.bigbluebutton.org/dev/api.html#getmeetings.
func (s *Server) GetMeetings(c *gin.Context) {
	instances, err := s.InstanceManager.ListInstances()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := &api.GetMeetingsResponse{
		ReturnCode: api.ReturnCodes().Success,
		Meetings:   make([]api.MeetingInfo, 0),
	}

	for _, instance := range instances {
		meetings, err := instance.GetMeetings()
		if err != nil {
			log.Error("An error occurred while retrieving meetings from instance", err)
			continue
		}

		response.Meetings = append(response.Meetings, meetings.Meetings...)
	}

	c.XML(http.StatusOK, response)
}

func missingMeetingIDParameter(c *gin.Context) {
	log.Error("Missing meetingID parameter")
	c.XML(http.StatusOK, api.CreateError(api.MessageKeys().ValidationError, api.Messages().EmptyMeetingID))
}

func missingRecordIDParameter(c *gin.Context) {
	log.Error("Missing recordID parameter")
	c.XML(http.StatusOK, api.CreateError(api.MessageKeys().MissingRecordIDParameter, api.Messages().MissingRecordIDParameter))
}

func (s *Server) retrieveBBBBInstanceFromMeetingID(meetingID string) (api.BigBlueButtonInstance, error) {
	host, err := s.Mapper.Get(MeetingMapKey(meetingID))
	if err != nil {
		return api.BigBlueButtonInstance{}, fmt.Errorf("SessionManager failed to retrieve session: %s", err.Error())
	}

	if host == "" {
		return api.BigBlueButtonInstance{}, errors.New("SessionManager failed to retrieve session host")
	}

	instance, err := s.InstanceManager.Get(host)
	if err != nil {
		return api.BigBlueButtonInstance{}, fmt.Errorf("manager failed to retrieve target instance for current request %s", err.Error())
	}

	return instance, nil
}

// Create handler find a server and create a meeting on balanced server.
func (s *Server) Create(c *gin.Context) {
	ctx := getAPIContext(c)
	instances, err := s.InstanceManager.List()
	if err != nil {
		log.Error("Manager failed to retrieve instance list", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if len(instances) == 0 {
		log.Error("InstanceManager does not retrieve any instances. Please check you add at least one Bigbluebutton instance")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	target, err := s.Balancer.Process(instances)
	if err != nil || target == "" {
		log.Error("Balancer failed to process current request", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	instance, err := s.InstanceManager.Get(target)
	if err != nil {
		log.Error("Manager failed to retrieve target instance for current request", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	apiResponse, err := instance.Create(ctx.Params)

	if err != nil {
		log.Error("An error occurred while creating remote session, instance returns a nil response", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	addErr := s.Mapper.Add(MeetingMapKey(apiResponse.MeetingID), instance.URL)
	if addErr != nil {
		log.Error("SessionManager failed to add new session", addErr)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.XML(http.StatusOK, apiResponse)
}

// Join handler join provided session. See https://docs.bigbluebutton.org/dev/api.html#join
func (s *Server) Join(c *gin.Context) {
	ctx := getAPIContext(c)
	meetingID, exists := c.GetQuery("meetingID")
	if !exists {
		missingMeetingIDParameter(c)
		return
	}

	redirect, redirectExists := c.GetQuery("redirect")

	instance, err := s.retrieveBBBBInstanceFromMeetingID(meetingID)
	if err != nil {
		log.Error(err)
		c.XML(http.StatusOK, api.CreateError(api.MessageKeys().NotFound, api.Messages().NotFound))
		return
	}

	if redirectExists && redirect == "false" {
		response, err := instance.Join(ctx.Params)
		if err != nil {
			log.Error("An error occurred while calling join instance api", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.XML(http.StatusOK, response)
	} else {
		redirectURL, err := instance.GetJoinRedirectURL(ctx.Params)
		if err != nil {
			log.Error("An error occurred while retrieving redirect URL on session join", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Redirect(http.StatusFound, redirectURL)
	}
}

// End handler end provided session. See https://docs.bigbluebutton.org/dev/api.html#end
func (s *Server) End(c *gin.Context) {
	endProcess := func() error {
		meetingID, _ := c.GetQuery("meetingID")
		removeErr := s.Mapper.Remove(MeetingMapKey(meetingID))
		if removeErr != nil {
			return fmt.Errorf("SessionManager failed to remove session %s: %s", meetingID, removeErr)
		}

		return nil
	}

	s.proxy(c, api.End, endProcess)
}

func (s *Server) proxy(c *gin.Context, action string, endProcess func() error) {
	ctx := getAPIContext(c)
	meetingID, exists := c.GetQuery("meetingID")
	if !exists {
		missingMeetingIDParameter(c)
		return
	}

	instance, err := s.retrieveBBBBInstanceFromMeetingID(meetingID)
	if err != nil {
		log.Error(err)
		c.XML(http.StatusOK, api.CreateError(api.MessageKeys().NotFound, api.Messages().NotFound))
		return
	}

	methodName := strings.Title(action)
	method := reflect.ValueOf(&instance).MethodByName(methodName)
	if method.IsNil() {
		log.Errorf("Failed to retrieve %s method on bigbluebutton instance", methodName)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	values := method.Call([]reflect.Value{reflect.ValueOf(ctx.Params)})
	response := values[0].Interface()
	responseErr := values[1].Interface()

	if responseErr != nil {
		log.Errorf("An error occurred while calling %s method on remote instance", action)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if endProcess != nil {
		err := endProcess()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	c.XML(http.StatusOK, response)
}

// IsMeetingRunning handler check if provided session is running. See https://docs.bigbluebutton.org/dev/api.html#ismeetingrunning
func (s *Server) IsMeetingRunning(c *gin.Context) {
	s.proxy(c, api.IsMeetingRunning, nil)
}

// GetMeetingInfo handler get information about provided session. See https://docs.bigbluebutton.org/dev/api.html#getmeetinginfo
func (s *Server) GetMeetingInfo(c *gin.Context) {
	s.proxy(c, api.GetMeetingInfo, nil)
}

// GetRecordings handler get recordings for provided session. See https://docs.bigbluebutton.org/dev/api.html#getrecordings
func (s *Server) GetRecordings(c *gin.Context) {
	ctx := getAPIContext(c)
	emptyRecordingsResponse := &api.GetRecordingsResponse{
		Response: api.Response{
			ReturnCode: api.ReturnCodes().Success,
			MessageKey: api.MessageKeys().NoRecordings,
			Message:    api.Messages().NoRecordings,
		},
	}

	instances, err := s.InstanceManager.ListInstances()
	if err != nil {
		log.Error("Manager failed to retrieve instances for getRecordings request", err)
		c.XML(http.StatusOK, emptyRecordingsResponse)
		return
	}

	response := &api.GetRecordingsResponse{
		Response: api.Response{
			ReturnCode: api.ReturnCodes().Success,
		},
		Recordings: []api.Recording{},
	}

	for _, instance := range instances {
		recordings, err := instance.GetRecordings(ctx.Params)
		if err != nil {
			log.Errorln(fmt.Sprintf("Instance %s failed to retrieve recordings.", instance.URL), err)
			continue
		}

		response.Recordings = append(response.Recordings, recordings.Recordings...)
	}

	if len(response.Recordings) == 0 {
		c.XML(http.StatusOK, emptyRecordingsResponse)
		return
	}

	c.XML(http.StatusOK, response)
}

// UpdateRecordings handler update recordings for provided record identifier. See https://docs.bigbluebutton.org/dev/api.html#updaterecordings
func (s *Server) UpdateRecordings(c *gin.Context) {
	ctx := getAPIContext(c)
	recordID, exists := c.GetQuery("recordID")
	if !exists {
		missingRecordIDParameter(c)
		return
	}

	instanceURL, err := s.Mapper.Get(RecordingMapKey(recordID))
	if err != nil {
		log.Errorln(fmt.Sprintf("Failed to retrieve instance for recording %s", recordID), err)
		c.XML(http.StatusOK, api.CreateError(api.MessageKeys().NotFound, api.Messages().NotFound))
		return
	}

	instance, err := s.InstanceManager.Get(instanceURL)
	if err != nil {
		log.Errorln(fmt.Sprintf("Failed to retrieve instance for instance url %s", recordID), err)
		c.XML(http.StatusOK, api.CreateError(api.MessageKeys().NotFound, api.Messages().NotFound))
		return
	}

	response, err := instance.UpdateRecordings(ctx.Params)
	if err != nil {
		log.Errorln(fmt.Sprintf("Instance %s failed to update recordings.", instance.URL), err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.XML(http.StatusOK, response)
}
