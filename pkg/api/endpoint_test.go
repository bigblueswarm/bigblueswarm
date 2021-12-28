package api

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAPIPath(t *testing.T) {
	assert.Equal(t, Path("action"), "/action")
}

func TestEndpointGroupLoad(t *testing.T) {
	router := gin.Default()

	group := EndpointGroup{
		Path: "/",
		Endpoints: []interface{}{
			Endpoint{
				Method:  "GET",
				Path:    "/",
				Handler: func(c *gin.Context) {},
			},
			EndpointGroup{
				Path: "/group",
				Endpoints: []interface{}{
					Endpoint{
						Method:  "GET",
						Path:    "/routes",
						Handler: func(c *gin.Context) {},
					},
				},
			},
		},
	}

	group.Load(router.Group(group.Path))

	assert.Equal(t, "/", router.BasePath())
	assert.Equal(t, 2, len(router.Routes()))
	assert.Equal(t, "/", router.Routes()[0].Path)
	assert.Equal(t, "/group/routes", router.Routes()[1].Path)
}
