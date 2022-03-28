package admin

import (
	"errors"
	"fmt"
	"testing"

	"github.com/SLedunois/b3lb/internal/test"
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
