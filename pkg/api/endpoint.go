package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Endpoint represent a server endpoint using a http method, a http path and a handler
type Endpoint struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
}

// EndpointGroup represent a group of endpoints
type EndpointGroup struct {
	Path      string
	Endpoints []interface{}
}

// Load endpoints from a list of endpoints
func (g *EndpointGroup) Load(group *gin.RouterGroup) {
	for _, endpoint := range g.Endpoints {
		if e, ok := endpoint.(Endpoint); ok {
			if e.Method == "" && e.Path == "" {
				group.Use(e.Handler)
			} else {
				group.Handle(e.Method, e.Path, e.Handler)
			}
		} else if e, ok := endpoint.(EndpointGroup); ok {
			e.Load(group.Group(e.Path))
		}
	}
}

// Path format action as a path
func Path(action string) string {
	return fmt.Sprintf("/%s", action)
}
