package admin

import (
	"errors"
	"fmt"
	"testing"

	"github.com/b3lb/test_utils/pkg/test"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestAddTenant(t *testing.T) {
	var tenant *Tenant
	tests := []test.Test{
		{
			Name: "adding a Tenant that does not contains a host in spec should return an error",
			Mock: func() {
				tenant = &Tenant{
					Kind:      "Tenant",
					Spec:      map[string]string{},
					Instances: []string{},
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			Name: "adding a Tenant should return an error if redis returns an error",
			Mock: func() {
				host := "localhost:8090"
				tenant = &Tenant{
					Kind: "Tenant",
					Spec: map[string]string{
						"host": host,
					},
					Instances: []string{
						"http://localhost/bigbluebutton",
					},
				}
				if out, err := yaml.Marshal(tenant); err == nil {
					mock := redisMock.ExpectSet(fmt.Sprintf("tenant:%s", host), string(out), 0)
					mock.SetVal("")
					mock.SetErr(errors.New("redis error"))
				} else {
					t.Error(err)
				}

			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			Name: "adding a Tenant should return nil if everything fine",
			Mock: func() {
				host := "localhost:8090"
				tenant = &Tenant{
					Kind: "Tenant",
					Spec: map[string]string{
						"host": host,
					},
					Instances: []string{
						"http://localhost/bigbluebutton",
					},
				}
				if out, err := yaml.Marshal(tenant); err == nil {
					redisMock.ExpectSet(fmt.Sprintf("tenant:%s", host), string(out), 0).SetVal("")
				} else {
					t.Error(err)
				}

			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Nil(t, err)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Mock()
			test.Validator(t, nil, tenantManager.AddTenant(tenant))
		})
	}
}

func TestListTenants(t *testing.T) {
	tests := []test.Test{
		{
			Name: "a redis error should return an error",
			Mock: func() {
				mock := redisMock.ExpectKeys("tenant:*")
				mock.SetErr(errors.New("redis error"))
				mock.SetVal([]string{})
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				res := value.([]TenantListObject)
				assert.NotNil(t, err)
				assert.Empty(t, res)
			},
		},
		{
			Name: "no error should return a valid list",
			Mock: func() {
				redisMock.ExpectKeys("tenant:*").SetVal([]string{"tenant:localhost"})
				tenant := &Tenant{
					Kind: "Tenant",
					Spec: map[string]string{
						"host": "localhost",
					},
					Instances: []string{
						"http://localhost/bigbluebutton",
					},
				}
				if expected, err := yaml.Marshal(tenant); err == nil {
					redisMock.ExpectGet("tenant:localhost").SetVal(string(expected))
				} else {
					t.Fatal(err)
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				res := value.([]TenantListObject)
				assert.Nil(t, err)
				assert.Equal(t, 1, len(res))
				assert.Equal(t, "localhost", res[0].Hostname)
				assert.Equal(t, 1, res[0].InstanceCount)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Mock()
			tenants, err := tenantManager.ListTenants()
			test.Validator(t, tenants, err)
		})
	}
}

func TestDeleteTenant(t *testing.T) {
	tests := []test.Test{
		{
			Name: "an error returned by redis should return an error",
			Mock: func() {
				mock := redisMock.ExpectDel("tenant:localhost")
				mock.SetVal(0)
				mock.SetErr(errors.New("redis error"))
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			Name: "a valid request should remove tenant from redis and return no error",
			Mock: func() {
				redisMock.ExpectDel("tenant:localhost").SetVal(1)
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Nil(t, err)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Mock()
			err := tenantManager.DeleteTenant("localhost")
			test.Validator(t, nil, err)
		})
	}
}

func TestGetTenant(t *testing.T) {
	tests := []test.Test{
		{
			Name: "an error returned by Redis should return the error",
			Mock: func() {
				mock := redisMock.ExpectGet("tenant:localhost")
				mock.SetVal("")
				mock.SetErr(errors.New("redis error"))
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			Name: "a not found tenant should return nil for tenant and error",
			Mock: func() {
				mock := redisMock.ExpectGet("tenant:localhost")
				mock.SetVal("")
				mock.SetErr(redis.Nil)
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Nil(t, err)
				assert.Nil(t, value)
			},
		},
		{
			Name: "a valid request should return a valid tenant",
			Mock: func() {
				tenant := Tenant{
					Kind: "Tenant",
					Spec: map[string]string{
						"host": "localhost",
					},
					Instances: []string{
						"http://localhost/bigbluebutton",
					},
				}

				if out, err := yaml.Marshal(tenant); err == nil {
					redisMock.ExpectGet("tenant:localhost").SetVal(string(out))
				} else {
					t.Error(err)
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				tenant := value.(*Tenant)
				assert.NotNil(t, tenant)
				assert.Equal(t, "localhost", tenant.Spec["host"])
				assert.Nil(t, err)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Mock()
			tenant, err := tenantManager.GetTenant("localhost")
			test.Validator(t, tenant, err)
		})
	}
}
