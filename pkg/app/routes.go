// Package app is the bigblueswarm core
package app

import (
	"net/http"

	"github.com/bigblueswarm/bigblueswarm/v2/pkg/admin"
	"github.com/bigblueswarm/bigblueswarm/v2/pkg/api"
)

func (s *Server) initRoutes() {
	adm := admin.CreateAdmin(s.InstanceManager, s.TenantManager, s.Balancer, s.Config)
	routes := append(*s.Routes(), *adm.Routes()...)
	for _, route := range routes {
		route.Load(s.Router.Group(route.Path))
	}
}

// Routes returns the server routes
func (s *Server) Routes() *[]api.EndpointGroup {
	return &[]api.EndpointGroup{
		{
			Path: api.Path(api.BigBlueButton),
			Endpoints: []interface{}{
				api.Endpoint{
					Method:  http.MethodGet,
					Handler: s.HealthCheck,
				},
				api.EndpointGroup{
					Path: api.Path(api.API),
					Endpoints: []interface{}{
						api.Endpoint{
							Method:  http.MethodGet,
							Handler: s.HealthCheck,
						},
						api.Endpoint{
							Handler: setLogger,
						},
						api.Endpoint{
							Handler: s.checkTenant,
						},
						api.Endpoint{
							Handler: s.ChecksumValidation,
						},
						api.Endpoint{
							Method:  http.MethodGet,
							Handler: s.Create,
							Path:    api.Path(api.Create),
						},
						api.Endpoint{
							Method:  http.MethodGet,
							Handler: s.GetMeetings,
							Path:    api.Path(api.GetMeetings),
						},
						api.Endpoint{
							Method:  http.MethodGet,
							Handler: s.Join,
							Path:    api.Path(api.Join),
						},
						api.Endpoint{
							Method:  http.MethodGet,
							Handler: s.End,
							Path:    api.Path(api.End),
						},
						api.Endpoint{
							Method:  http.MethodGet,
							Handler: s.IsMeetingRunning,
							Path:    api.Path(api.IsMeetingRunning),
						},
						api.Endpoint{
							Method:  http.MethodGet,
							Handler: s.GetMeetingInfo,
							Path:    api.Path(api.GetMeetingInfo),
						},
						api.Endpoint{
							Method:  http.MethodGet,
							Handler: s.GetRecordings,
							Path:    api.Path(api.GetRecordings),
						},
						api.Endpoint{
							Method:  http.MethodGet,
							Handler: s.UpdateRecordings,
							Path:    api.Path(api.UpdateRecordings),
						},
						api.Endpoint{
							Method:  http.MethodGet,
							Handler: s.DeleteRecordings,
							Path:    api.Path(api.DeleteRecordings),
						},
						api.Endpoint{
							Method:  http.MethodGet,
							Handler: s.PublishRecordings,
							Path:    api.Path(api.PublishRecordings),
						},
						api.Endpoint{
							Method:  http.MethodGet,
							Handler: s.GetRecordingsTextTracks,
							Path:    api.Path(api.GetRecordingsTextTracks),
						},
						api.Endpoint{
							Method:  http.MethodPost,
							Handler: s.PutRecordingTextTrack,
							Path:    api.Path(api.PutRecordingTextTrack),
						},
					},
				},
			},
		},
	}
}
