package app

import (
	"b3lb/pkg/admin"
	"b3lb/pkg/api"
	"net/http"
)

func (s *Server) initRoutes() {
	adm := admin.CreateAdmin(s.InstanceManager, &s.Config.Admin)
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
					},
				},
			},
		},
	}
}
