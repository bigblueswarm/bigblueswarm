// Package app is the bigblueswarm core
package app

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/bigblueswarm/bigblueswarm/v3/pkg/admin"
	"github.com/bigblueswarm/bigblueswarm/v3/pkg/api"
	"github.com/bigblueswarm/bigblueswarm/v3/pkg/utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// HealthCheck handler returns an health check response
func (s *Server) HealthCheck(c *gin.Context) {
	c.XML(http.StatusOK, api.CreateHealthCheck())
}

// checkTenant check if tenant exists. Otherwise, it returns an error
func (s *Server) checkTenant(c *gin.Context) {
	hostname := utils.GetHost(c)
	logger := getLogger(c)
	tenant, err := s.TenantManager.GetTenant(hostname)
	if err != nil {
		logger.Errorln("failed to retrieve tenant", err)
		c.XML(http.StatusInternalServerError, getTenantError())
		c.Abort()
		return
	}

	logger.addField("tenant", utils.GetHost(c))
	if tenant == nil {
		logger.Warn("tenant manager does not found current hostname")
		c.XML(http.StatusForbidden, tenantNotFoundError())
		c.Abort()
		return
	}

	c.Next()
}

// GetMeetings handler returns the getMeetings API. See https://docs.bigbluebutton.org/dev/api.html#getmeetings.
func (s *Server) GetMeetings(c *gin.Context) {
	logger := getLogger(c)
	instances, err := s.InstanceManager.ListInstances()
	if err != nil {
		logger.Errorln("failed to retrieving instances", err)
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
			logger.dup().addField("instance", instance.URL).Errorln("An error occurred while retrieving meetings from instance", err)
			continue
		}

		response.Meetings = append(response.Meetings, meetings.Meetings...)
	}

	c.XML(http.StatusOK, response)
}

func missingMeetingIDParameter(c *gin.Context) {
	c.XML(http.StatusOK, api.CreateError(api.MessageKeys().ValidationError, api.Messages().EmptyMeetingID))
}

func missingRecordIDParameter(c *gin.Context) {
	c.XML(http.StatusOK, api.CreateError(api.MessageKeys().MissingRecordIDParameter, api.Messages().MissingRecordIDParameter))
}

func (s *Server) retrieveBBBBInstanceFromKey(key string) (api.BigBlueButtonInstance, error) {
	host, err := s.Mapper.Get(key)
	if err != nil {
		return api.BigBlueButtonInstance{}, fmt.Errorf("mapper failed to retrieve session: %s", err.Error())
	}

	if host == "" {
		return api.BigBlueButtonInstance{}, errors.New("mapper failed to retrieve session host")
	}

	instance, err := s.InstanceManager.Get(host)
	if err != nil {
		return api.BigBlueButtonInstance{}, fmt.Errorf("manager failed to retrieve target instance for current request %s", err.Error())
	}

	return instance, nil
}

func (s *Server) canTenantJoinMeeting(logger *RequestLogger, tenant *admin.Tenant) (int, *api.Error) {
	if tenant.HasUserPool() {
		logger := logger.setFields(log.Fields{
			"tenant":        tenant.Spec.Host,
			"meetings_pool": tenant.Spec.MeetingsPool,
			"users_pool":    tenant.Spec.UserPool,
		})

		canCreate, err := s.isTenantLowerThanUserPool(tenant)
		if err != nil {
			logger.Errorln("unable to check if tenant can join meeting", err)
			return http.StatusInternalServerError, serverError("BigBlueSwarm failed to check if your tenant reached the user pool limit.")
		}

		if !canCreate {
			log.Info("tenant raise the user pool limit and can't join meeting")
			return http.StatusForbidden, userPoolReachedError()
		}
	}

	return http.StatusOK, nil
}

func (s *Server) canTenantCreateMeeting(logger *RequestLogger, tenant *admin.Tenant) (int, *api.Error) {
	if tenant.HasMeetingPool() {
		logger := logger.setFields(log.Fields{
			"tenant":        tenant.Spec.Host,
			"meetings_pool": tenant.Spec.MeetingsPool,
			"users_pool":    tenant.Spec.UserPool,
		})

		canCreate, err := s.isTenantLowerThanMeetingPool(tenant)
		if err != nil {
			logger.Errorln("unable to check if tenant can create meeting", err)
			return http.StatusInternalServerError, serverError("BigBlueSwarm failed to check if your tenant reached the meeting pool limit.")
		}

		if !canCreate {
			logger.Info("tenant raise the meetings pool limit and can't create a new meeting")
			return http.StatusForbidden, meetingPoolReachedError()
		}
	}

	return http.StatusOK, nil
}

// Create handler find a server and create a meeting on balanced server.
func (s *Server) Create(c *gin.Context) {
	ctx := getAPIContext(c)
	tenant, err := s.TenantManager.GetTenant(utils.GetHost(c))
	logger := getLogger(c)

	if err != nil {
		logger.Errorln("manager failed to retrieve tenant", err)
		c.XML(http.StatusInternalServerError, getTenantError())
		return
	}

	logger.setFields(log.Fields{
		"tenant": tenant.Spec.Host,
		"action": ctx.Action,
		"params": ctx.Params,
	})

	if status, err := s.canTenantCreateMeeting(logger.dup(), tenant); status != http.StatusOK {
		c.XML(status, err)
		return
	}

	if len(tenant.Instances) == 0 {
		logger.Info("tenant does not have a configured instance list. Getting all instances")
		instances, err := s.InstanceManager.List()
		if err != nil {
			logger.Errorln("instance manager failed to retrieve instances", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		tenant.Instances = instances
	}

	ctx.SetTenantMetadata(tenant.Spec.Host)

	target, err := s.Balancer.Process(tenant.Instances)
	if err != nil || target == "" {
		logger.Errorln("balancer failed to process current request", err)
		c.XML(http.StatusInternalServerError, noInstanceFoundError())
		return
	}

	instance, err := s.InstanceManager.Get(target)
	if err != nil {
		logger.Errorln("manager failed to retrieve target instance for current request", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	apiResponse, err := instance.Create(ctx.Params)

	if err != nil {
		logger.Errorln("an error occurred while creating remote session, instance returns a nil response", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	addErr := s.Mapper.Add(MeetingMapKey(apiResponse.MeetingID), instance.URL)
	if addErr != nil {
		logger.Errorln("mapper failed to add new session", addErr)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.XML(http.StatusOK, apiResponse)
}

// Join handler join provided session. See https://docs.bigbluebutton.org/dev/api.html#join
func (s *Server) Join(c *gin.Context) {
	ctx := getAPIContext(c)
	logger := getLogger(c)
	tenant, err := s.TenantManager.GetTenant(utils.GetHost(c))

	if err != nil {
		logger.Errorln("failed to retrieve tenant from host", err)
		c.XML(http.StatusInternalServerError, getTenantError())
		return
	}

	logger.setFields(log.Fields{
		"tenant": tenant.Spec.Host,
		"action": ctx.Action,
		"params": ctx.Params,
	})

	if status, err := s.canTenantJoinMeeting(logger.dup(), tenant); status != http.StatusOK {
		c.XML(status, err)
		return
	}

	meetingID, exists := c.GetQuery("meetingID")
	if !exists {
		logger.Warn("meeting id parameter missing")
		missingMeetingIDParameter(c)
		return
	}

	redirect, redirectExists := c.GetQuery("redirect")

	instance, err := s.retrieveBBBBInstanceFromKey(MeetingMapKey(meetingID))
	if err != nil {
		logger.Error(err)
		c.XML(http.StatusOK, api.CreateError(api.MessageKeys().NotFound, api.Messages().NotFound))
		return
	}

	if redirectExists && redirect == "false" {
		response, err := instance.Join(ctx.Params)
		if err != nil {
			logger.Errorln("An error occurred while calling join instance api", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.XML(http.StatusOK, response)
	} else {
		redirectURL, err := instance.GetJoinRedirectURL(ctx.Params)
		if err != nil {
			logger.Errorln("An error occurred while retrieving redirect URL on session join", err)
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
			return fmt.Errorf("mapper failed to remove session %s: %s", meetingID, removeErr)
		}

		return nil
	}

	s.proxy(c, api.End, endProcess)
}

func errorMessage(action string) interface{} {
	switch action {
	case api.GetRecordingsTextTracks:
		return api.CreateJSONError(api.MessageKeys().NoRecordings, api.Messages().RecordingTextTrackNotFound)
	default:
		return api.CreateError(api.MessageKeys().NotFound, api.Messages().RecordingNotFound)
	}
}

func ginMethod(action string, c *gin.Context) func(code int, obj interface{}) {
	switch action {
	case api.GetRecordingsTextTracks:
		return c.JSON
	default:
		return c.XML
	}
}

func (s *Server) proxyRecordings(c *gin.Context, action string, endProcess func(response interface{}) error) {
	ctx := getAPIContext(c)
	logger := getLogger(c)
	recordID, exists := c.GetQuery("recordID")
	if !exists {
		logger.Warn("record id parameter missing")
		missingRecordIDParameter(c)
		return
	}

	logger.addField("record", recordID)

	instance, err := s.retrieveBBBBInstanceFromKey(RecordingMapKey(recordID))
	if err != nil {
		logger.Errorln("failed to retrieve instance", err)
		ginMethod(action, c)(http.StatusOK, errorMessage(action))
		return
	}

	logger.addField("instance", instance.URL)
	response, mErr := callInstanceMethod(ctx, instance, action)
	if mErr != nil {
		logger.Error(err)
		c.XML(http.StatusInternalServerError, serverError("BigBlueSwarm failed to call remote instance method"))
		return
	}

	if endProcess != nil {
		err := endProcess(response)
		if err != nil {
			logger.Error(err)
			c.XML(http.StatusInternalServerError, serverError("BigBlueSwarm failed to process recording method"))
			return
		}
	}

	ginMethod(action, c)(http.StatusOK, response)
}

func callInstanceMethod(ctx *api.Checksum, instance api.BigBlueButtonInstance, action string) (interface{}, interface{}) {
	methodName := strings.Title(action)
	value := reflect.ValueOf(&instance)
	if value.IsNil() {
		return nil, errors.New("failed to execute reflect on instance")
	}

	method := value.MethodByName(methodName)
	if method.IsNil() {
		return nil, fmt.Errorf("failed to retrieve %s method on bigbluebutton instance", methodName)
	}

	values := method.Call([]reflect.Value{reflect.ValueOf(ctx.Params)})
	return values[0].Interface(), values[1].Interface()
}

func (s *Server) proxy(c *gin.Context, action string, endProcess func() error) {
	ctx := getAPIContext(c)
	logger := getLogger(c)
	meetingID, exists := c.GetQuery("meetingID")
	if !exists {
		logger.Error("missing meeting id parameter")
		missingMeetingIDParameter(c)
		return
	}

	logger.addField("meeting_id", meetingID)

	instance, err := s.retrieveBBBBInstanceFromKey(MeetingMapKey(meetingID))
	if err != nil {
		logger.Error(err)
		ginMethod(action, c)(http.StatusOK, api.CreateError(api.MessageKeys().NotFound, api.Messages().NotFound))
		return
	}

	response, mErr := callInstanceMethod(ctx, instance, action)
	if mErr != nil {
		logger.Error(err)
		c.XML(http.StatusInternalServerError, serverError("BigBlueSwarm failed to process api call"))
		return
	}

	if endProcess != nil {
		err := endProcess()
		if err != nil {
			logger.Error(err)
			c.XML(http.StatusInternalServerError, serverError("BigBlueSwarm failed to end api process"))
			return
		}
	}

	ginMethod(action, c)(http.StatusOK, response)
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
	logger := getLogger(c).addField("action", ctx.Action)
	emptyRecordingsResponse := &api.GetRecordingsResponse{
		Response: api.Response{
			ReturnCode: api.ReturnCodes().Success,
			MessageKey: api.MessageKeys().NoRecordings,
			Message:    api.Messages().NoRecordings,
		},
	}

	instances, err := s.InstanceManager.ListInstances()
	if err != nil {
		logger.Errorln("manager failed to retrieve instances for getRecordings request", err)
		c.XML(http.StatusInternalServerError, serverError("BigBlueSwarm failed to process GetRecordings method"))
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
			logger.dup().addField("instance", instance.URL).Errorln("instance failed to retrieve recordings.", err)
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
	s.proxyRecordings(c, api.UpdateRecordings, nil)
}

// DeleteRecordings handler delete a single recording for provided record identifier. See https://docs.bigbluebutton.org/dev/api.html#deleterecordings
func (s *Server) DeleteRecordings(c *gin.Context) {
	endProcess := func(response interface{}) error {
		if deletion, ok := response.(*api.DeleteRecordingsResponse); ok && deletion.Deleted {
			recordID, _ := c.GetQuery("recordID")
			return s.Mapper.Remove(RecordingMapKey(recordID))
		}

		return nil
	}

	s.proxyRecordings(c, api.DeleteRecordings, endProcess)
}

// PublishRecordings handler publish a single recording for provided record identifier. See https://docs.bigbluebutton.org/dev/api.html#publishrecordings
func (s *Server) PublishRecordings(c *gin.Context) {
	s.proxyRecordings(c, api.PublishRecordings, nil)
}

// GetRecordingsTextTracks handler retrieve  list of the caption/subtitle tracks for a recording. See https://docs.bigbluebutton.org/dev/api.html#getrecordingstexttracks
func (s *Server) GetRecordingsTextTracks(c *gin.Context) {
	s.proxyRecordings(c, api.GetRecordingsTextTracks, nil)
}

// PutRecordingTextTrack handler redirect to the right bigbluebutton instance
func (s *Server) PutRecordingTextTrack(c *gin.Context) {
	ctx := getAPIContext(c)
	logger := getLogger(c)
	recordID, exists := c.GetQuery("recordID")
	if !exists {
		logger.Warn("missing record id parameter")
		c.AbortWithStatusJSON(http.StatusOK, api.CreateJSONError(api.MessageKeys().ParamError, api.Messages().MissingRecordIDParameter))
		return
	}

	logger.addField("record_id", recordID)
	instance, err := s.retrieveBBBBInstanceFromKey(RecordingMapKey(recordID))
	if err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusOK, api.CreateJSONError(api.MessageKeys().NoRecordings, api.Messages().RecordingTextTrackNotFound))
		return
	}

	instance.Redirect(c, api.PutRecordingTextTrack, ctx.Params)
}
