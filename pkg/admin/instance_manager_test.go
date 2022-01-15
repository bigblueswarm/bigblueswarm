package admin

import (
	"errors"
	"testing"

	"github.com/SLedunois/b3lb/internal/test"
	"github.com/SLedunois/b3lb/pkg/api"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

const url string = "https://bbb_test.com"

func TestInstanceManagerExists(t *testing.T) {
	tests := []test.Test{
		{
			Name: "Existing value should return true and no error",
			Mock: func() {
				mock := redisMock.ExpectHExists(B3LBInstances, url)
				mock.SetErr(nil)
				mock.SetVal(true)
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.True(t, value.(bool))
				assert.Nil(t, err)
			},
		},
		{
			Name: "Non existing value should return true and no error",
			Mock: func() {
				mock := redisMock.ExpectHExists(B3LBInstances, url)
				mock.SetErr(nil)
				mock.SetVal(false)
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.False(t, value.(bool))
				assert.Nil(t, err)
			},
		},
		{
			Name: "Returning redis.Nil should return false and no error",
			Mock: func() {
				mock := redisMock.ExpectHExists(B3LBInstances, url)
				mock.SetErr(redis.Nil)
				mock.SetVal(false)
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.False(t, value.(bool))
				assert.Nil(t, err)
			},
		},
		{
			Name: "Returning an error should return the error",
			Mock: func() {
				mock := redisMock.ExpectHExists(B3LBInstances, url)
				mock.SetErr(errors.New("test error"))
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.NotNil(t, err)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Mock()
			exists, err := instanceManager.Exists(api.BigBlueButtonInstance{URL: url})
			test.Validator(t, exists, err)
		})
	}
}

func TestInstanceManagerList(t *testing.T) {
	tests := []test.Test{
		{
			Name: "Returning a list of keys should return the list and no error",
			Mock: func() {
				redisMock.ExpectHKeys(B3LBInstances).SetVal([]string{url})
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, value.([]string)[0], url)
				assert.Nil(t, err)
			},
		},
		{
			Name: "Returning redis.Nil should return an empty list and no error",
			Mock: func() {
				mock := redisMock.ExpectHKeys(B3LBInstances)
				mock.SetErr(redis.Nil)
				mock.SetVal([]string{})
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, len(value.([]string)), 0)
				assert.Nil(t, err)
			},
		},
		{
			Name: "Returning an error should return the error",
			Mock: func() {
				redisMock.ExpectHKeys(B3LBInstances).SetErr(errors.New("test error"))
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.NotNil(t, err)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Mock()
			keys, err := instanceManager.List()
			test.Validator(t, keys, err)
		})
	}
}

func TestInstanceManagerAdd(t *testing.T) {
	tests := []test.Test{
		{
			Name: "Adding a new instance should return no error",
			Mock: func() {
				redisMock.ExpectHSet(B3LBInstances, url, "secret")
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			Name: "Throwing an error when adding an instance should return the error",
			Mock: func() {
				redisMock.ExpectHSet(B3LBInstances, url, "secret").SetErr(errors.New("test error"))
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.NotNil(t, err)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Mock()
			err := instanceManager.Add(api.BigBlueButtonInstance{URL: url, Secret: "secret"})
			test.Validator(t, nil, err)
		})
	}
}

func TestInstanceManagerRemove(t *testing.T) {
	tests := []test.Test{
		{
			Name: "Removing an instance should return no error",
			Mock: func() {
				mock := redisMock.ExpectHDel(B3LBInstances, url)
				mock.SetErr(nil)
				mock.SetVal(1)
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Nil(t, err)
			},
		},
		{
			Name: "Throwing an error when removing an instance should return the error",
			Mock: func() {
				redisMock.ExpectHDel(B3LBInstances, url).SetErr(errors.New("test error"))
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.NotNil(t, err)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Mock()
			err := instanceManager.Remove(url)
			test.Validator(t, nil, err)
		})
	}
}

func TestInstanceManagerGet(t *testing.T) {
	secret := "secret"
	tests := []test.Test{
		{
			Name: "Getting an instance should return the instance and no error",
			Mock: func() {
				mock := redisMock.ExpectHGet(B3LBInstances, url)
				mock.SetErr(nil)
				mock.SetVal(secret)
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				instance := value.(api.BigBlueButtonInstance)
				assert.Equal(t, instance.URL, url)
				assert.Equal(t, instance.Secret, secret)
				assert.Nil(t, err)
			},
		},
		{
			Name: "Throwing an error when getting an instance should return the error",
			Mock: func() {
				mock := redisMock.ExpectHGet(B3LBInstances, url)
				mock.SetErr(errors.New("test error"))
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.NotNil(t, err)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Mock()
			instance, err := instanceManager.Get(url)
			test.Validator(t, instance, err)
		})
	}
}

func TestInstanceManagerListInstances(t *testing.T) {
	tests := []test.Test{
		{
			Name: "An empty map should return an empty list",
			Mock: func() {
				redisMock.ExpectHGetAll(B3LBInstances).SetVal(map[string]string{})
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				instances := value.([]api.BigBlueButtonInstance)
				assert.Nil(t, err)
				assert.Equal(t, len(instances), 0)
			},
		},
		{
			Name: "A map contanaing one instance should return a list with one instance",
			Mock: func() {
				instances := map[string]string{
					"http://localhost/bigbluebutton": "secret",
				}
				redisMock.ExpectHGetAll(B3LBInstances).SetVal(instances)
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				instances := value.([]api.BigBlueButtonInstance)
				assert.Nil(t, err)
				assert.Equal(t, len(instances), 1)
				assert.Equal(t, instances[0].URL, "http://localhost/bigbluebutton")
				assert.Equal(t, instances[0].Secret, "secret")
			},
		},
		{
			Name: "Redis returning an error should return an error and an empty list",
			Mock: func() {
				redisMock.ExpectHGetAll(B3LBInstances).SetErr(errors.New("redis error"))
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				instances := value.([]api.BigBlueButtonInstance)
				assert.NotNil(t, err)
				assert.Equal(t, len(instances), 0)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Mock()
			instances, err := instanceManager.ListInstances()
			test.Validator(t, instances, err)
		})
	}
}
