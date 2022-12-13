package app

import (
	"errors"
	"testing"

	"github.com/bigblueswarm/bigblueswarm/v2/pkg/admin"
	"github.com/bigblueswarm/bigblueswarm/v2/pkg/balancer"
	"github.com/bigblueswarm/bigblueswarm/v2/pkg/config"
	"github.com/bigblueswarm/test_utils/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestIsTenantLowerThanMeetingPool(t *testing.T) {
	mp := int64(0)
	tenant := &admin.Tenant{
		Spec: &admin.TenantSpec{
			MeetingsPool: &mp,
		},
	}
	server := NewServer(&config.Config{})
	server.Balancer = &balancer.Mock{}
	tests := []test.Test{
		{
			Name: "an error returned by balancer should be returned",
			Mock: func() {
				balancer.BalancerGetCurrentStateFunc = func(measurement, field string) (int64, error) {
					return 0, errors.New("balancer error")
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			Name: "a state lower than the constraint should return true",
			Mock: func() {
				constraint := int64(100)
				tenant.Spec.MeetingsPool = &constraint
				balancer.BalancerGetCurrentStateFunc = func(measurement, field string) (int64, error) {
					return 99, nil
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Nil(t, err)
				assert.Equal(t, true, value.(bool))
			},
		},
		{
			Name: "a state hidher than the constraint should return true",
			Mock: func() {
				constraint := int64(90)
				tenant.Spec.MeetingsPool = &constraint
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Nil(t, err)
				assert.Equal(t, false, value.(bool))
			},
		},
	}

	for _, test := range tests {
		test.Mock()
		can, err := server.isTenantLowerThanMeetingPool(tenant)
		test.Validator(t, can, err)
	}
}

func TestIsTenantLowerThanUserPool(t *testing.T) {
	up := int64(0)
	tenant := &admin.Tenant{
		Spec: &admin.TenantSpec{
			UserPool: &up,
		},
	}
	server := NewServer(&config.Config{})
	server.Balancer = &balancer.Mock{}
	tests := []test.Test{
		{
			Name: "an error returned by balancer should be returned",
			Mock: func() {
				balancer.BalancerGetCurrentStateFunc = func(measurement, field string) (int64, error) {
					return 0, errors.New("balancer error")
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			Name: "a state lower than the constraint should return true",
			Mock: func() {
				constraint := int64(100)
				tenant.Spec.UserPool = &constraint
				balancer.BalancerGetCurrentStateFunc = func(measurement, field string) (int64, error) {
					return 99, nil
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Nil(t, err)
				assert.Equal(t, true, value.(bool))
			},
		},
		{
			Name: "a state hidher than the constraint should return true",
			Mock: func() {
				constraint := int64(90)
				tenant.Spec.UserPool = &constraint
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Nil(t, err)
				assert.Equal(t, false, value.(bool))
			},
		},
	}

	for _, test := range tests {
		test.Mock()
		can, err := server.isTenantLowerThanUserPool(tenant)
		test.Validator(t, can, err)
	}
}
