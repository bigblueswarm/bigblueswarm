package admin

import (
	"net/http"

	"github.com/SLedunois/b3lb/pkg/api"
)

// Routes returns admin routes
func (a *Admin) Routes() *[]api.EndpointGroup {
	return &[]api.EndpointGroup{
		{
			Path: "/admin",
			Endpoints: []interface{}{
				api.Endpoint{
					Handler: a.APIKeyValidation,
				},
				api.EndpointGroup{
					Path: "/api",
					Endpoints: []interface{}{
						api.EndpointGroup{
							Path: "/instances",
							Endpoints: []interface{}{
								api.Endpoint{
									Method:  http.MethodGet,
									Handler: a.ListInstances,
								},
								api.Endpoint{
									Method:  http.MethodPost,
									Handler: a.SetInstances,
								},
							},
						},
						api.EndpointGroup{
							Path: "/tenants",
							Endpoints: []interface{}{
								api.Endpoint{
									Method:  http.MethodPost,
									Handler: a.CreateTenant,
								},
							},
						},
						api.Endpoint{
							Path:    "/cluster",
							Method:  http.MethodGet,
							Handler: a.ClusterStatus,
						},
					},
				},
			},
		},
	}
}
