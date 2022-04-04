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
									Method:  http.MethodGet,
									Handler: a.ListTenants,
								},
								api.Endpoint{
									Method:  http.MethodPost,
									Handler: a.CreateTenant,
								},
								api.EndpointGroup{
									Path: "/:hostname",
									Endpoints: []interface{}{
										api.Endpoint{
											Method:  http.MethodDelete,
											Handler: a.DeleteTenant,
										},
										api.Endpoint{
											Method:  http.MethodGet,
											Handler: a.GetTenant,
										},
									},
								},
							},
						},
						api.Endpoint{
							Path:    "/cluster",
							Method:  http.MethodGet,
							Handler: a.ClusterStatus,
						},
						api.Endpoint{
							Path:    "/configurations",
							Method:  http.MethodGet,
							Handler: a.GetConfiguration,
						},
					},
				},
			},
		},
	}
}
