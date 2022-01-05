package app

import (
	"b3lb/pkg/api"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// HealthCheck handler returns an health check response
func (s *Server) HealthCheck(c *gin.Context) {
	c.XML(http.StatusOK, api.CreateHealthCheck())
}

// GetMeetings handler returns the getMeetings API. See https://docs.bigbluebutton.org/dev/api.html#getmeetings.
func (s *Server) GetMeetings(c *gin.Context) {
	c.String(http.StatusOK, c.FullPath())
}

func missingMeetingIDParameter(c *gin.Context) {
	log.Error("Missing meetingID parameter")
	c.XML(http.StatusOK, api.CreateError(api.MessageKeys().ValidationError, api.Messages().EmptyMeetingID))
}

func (s *Server) retrieveBBBBInstanceFromMeetingID(meetingID string) (api.BigBlueButtonInstance, error) {
	host, err := s.SessionManager.Get(meetingID)
	if err != nil {
		return api.BigBlueButtonInstance{}, fmt.Errorf("SessionManager failed to retrieve session: %s", err.Error())
	}

	if host == "" {
		return api.BigBlueButtonInstance{}, errors.New("SessionManager failed to retrieve session host")
	}

	instance, err := s.InstanceManager.Get(host)
	if err != nil {
		return api.BigBlueButtonInstance{}, fmt.Errorf("Manager failed to retrieve target instance for current request %s", err.Error())
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

	apiResponse := instance.Create(ctx.Params)

	if apiResponse == nil {
		log.Error("An error occurred while creating remote session, instance returns a nil response")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	addErr := s.SessionManager.Add(apiResponse.MeetingID, instance.URL)
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

	instance, err := s.retrieveBBBBInstanceFromMeetingID(meetingID)
	if err != nil {
		log.Error(err)
		c.XML(http.StatusOK, api.CreateError(api.MessageKeys().NotFound, api.Messages().NotFound))
		return
	}

	redirectURL, err := instance.GetJoinRedirectURL(ctx.Params)
	if err != nil {
		log.Error("An error occurred while retrieving redirect URL on session join", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, redirectURL)
}

// End handler end provided session. See https://docs.bigbluebutton.org/dev/api.html#end
func (s *Server) End(c *gin.Context) {
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

	endResponse := instance.End(ctx.Params)
	if endResponse == nil {
		log.Error("An error occurred while ending remote session")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	removeErr := s.SessionManager.Remove(meetingID)
	if removeErr != nil {
		log.Error("SessionManager failed to remove session", removeErr)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.XML(http.StatusOK, endResponse)
}

func (s *Server) proxy(c *gin.Context, action string) {
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

	response := method.Call([]reflect.Value{reflect.ValueOf(ctx.Params)})[0].Interface()

	if response == nil {
		log.Errorf("An error occurred while calling %s method on remote instance", action)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.XML(http.StatusOK, response)
}

// IsMeetingRunning handler check if provided session is running. See https://docs.bigbluebutton.org/dev/api.html#ismeetingrunning
func (s *Server) IsMeetingRunning(c *gin.Context) {
	s.proxy(c, api.IsMeetingRunning)
}

// GetMeetingInfo handler get information about provided session. See https://docs.bigbluebutton.org/dev/api.html#getmeetinginfo
func (s *Server) GetMeetingInfo(c *gin.Context) {
	s.proxy(c, api.GetMeetingInfo)
}
