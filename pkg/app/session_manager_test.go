package app

import (
	"errors"
	"testing"

	"github.com/SLedunois/b3lb/internal/test"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

const sessionID string = "session_id"
const host string = "http://localhost/bigbluebutton"

func TestAdd(t *testing.T) {
	tests := []test.Test{
		{
			Name: "Add should not return and error if the value is added",
			Mock: func() {
				mock := redisMock.ExpectSet(sessionID, host, 0)
				mock.SetErr(redis.Nil)
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Nil(t, err)
			},
		},
		{
			Name: "Add should return an error if redis throws an error",
			Mock: func() {
				mock := redisMock.ExpectSet(sessionID, host, 0)
				mock.SetErr(errors.New("redis error"))
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.NotNil(t, err)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Mock()
			err := sessionManager.Add(sessionID, host)
			test.Validator(t, nil, err)
		})
	}
}

func TestGet(t *testing.T) {
	tests := []test.Test{
		{
			Name: "Get should return the host and no error if the session is found",
			Mock: func() {
				redisMock.ExpectGet(sessionID).SetVal(host)
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Nil(t, err)
				assert.Equal(t, host, value.(string))
			},
		},
		{
			Name: "Get should return an error if redis throws an error",
			Mock: func() {
				mock := redisMock.ExpectGet(sessionID)
				mock.SetErr(errors.New("redis error"))
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			Name: "Get should return an empty string if the session is not found",
			Mock: func() {
				redisMock.ExpectGet(sessionID).SetVal("")
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Nil(t, err)
				assert.Equal(t, "", value.(string))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Mock()
			value, err := sessionManager.Get(sessionID)
			test.Validator(t, value, err)
		})
	}
}

func TestRemove(t *testing.T) {
	tests := []test.Test{
		{
			Name: "Remove should return nil if the session is removed",
			Mock: func() {
				redisMock.ExpectDel(sessionID).SetErr(redis.Nil)
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Nil(t, err)
			},
		},
		{
			Name: "Remove should return an error if redis throws an error",
			Mock: func() {
				redisMock.ExpectDel(sessionID).SetErr(errors.New("redis error"))
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.NotNil(t, err)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Mock()
			err := sessionManager.Remove(sessionID)
			test.Validator(t, nil, err)
		})
	}
}
