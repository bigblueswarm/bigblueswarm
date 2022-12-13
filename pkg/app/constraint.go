// Package app is the bigblueswarm core
package app

import (
	"fmt"

	"github.com/bigblueswarm/bigblueswarm/v2/pkg/admin"
)

func (s *Server) canTenantCreateMeeting(t *admin.Tenant) (bool, error) {
	measurement := fmt.Sprintf("%s:bigbluebutton_meetings", t.Spec.Host)
	field := "active_meetings"
	status, err := s.Balancer.GetCurrentState(measurement, field)
	if err != nil {
		return false, fmt.Errorf("failed to check tenant state for tenant %s: %s", t.Spec.Host, err)
	}

	return status < *t.Spec.MeetingsPool, nil
}
