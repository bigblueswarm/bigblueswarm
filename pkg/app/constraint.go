// Package app is the bigblueswarm core
package app

import (
	"fmt"

	"github.com/bigblueswarm/bigblueswarm/v3/pkg/admin"
	"github.com/bigblueswarm/bigblueswarm/v3/pkg/balancer"
)

func isPoolReached(b balancer.Balancer, m string, f string, p int64) (bool, error) {
	state, err := b.GetCurrentState(m, f)
	if err != nil {
		return false, err
	}

	return state < p, nil
}

func (s *Server) isTenantLowerThanMeetingPool(t *admin.Tenant) (bool, error) {
	measurement := fmt.Sprintf("bigbluebutton:%s", t.Spec.Host)
	field := "meetings"
	reached, err := isPoolReached(s.Balancer, measurement, field, *t.Spec.MeetingsPool)
	if err != nil {
		return false, fmt.Errorf("failed to check tenant state for tenant %s: %s", t.Spec.Host, err)
	}

	return reached, nil
}

func (s *Server) isTenantLowerThanUserPool(t *admin.Tenant) (bool, error) {
	measurement := fmt.Sprintf("bigbluebutton:%s", t.Spec.Host)
	field := "participants"
	reached, err := isPoolReached(s.Balancer, measurement, field, *t.Spec.UserPool)
	if err != nil {
		return false, fmt.Errorf("failed to check tenant state for tenant %s: %s", t.Spec.Host, err)
	}

	return reached, nil
}
