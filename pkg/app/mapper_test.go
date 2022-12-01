package app

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bigblueswarm/test_utils/pkg/test"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

const id string = "20d93ede-ea0e-4df5-931f-5a858b994a28"
const host string = "http://localhost/bigbluebutton"

func TestMeetingMapKey(t *testing.T) {
	assert.Equal(t, fmt.Sprintf("meeting:%s", id), MeetingMapKey(id))
}

func TestRecordingMapKey(t *testing.T) {
	assert.Equal(t, fmt.Sprintf("recording:%s", id), RecordingMapKey(id))
}

func TestAdd(t *testing.T) {
	tests := []test.Test{
		{
			Name: "Add should not return and error if the value is added",
			Mock: func() {
				mock := redisMock.ExpectSet(MeetingMapKey(id), host, 0)
				mock.SetErr(redis.Nil)
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Nil(t, err)
			},
		},
		{
			Name: "Add should return an error if redis throws an error",
			Mock: func() {
				mock := redisMock.ExpectSet(MeetingMapKey(id), host, 0)
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
			err := mapper.Add(MeetingMapKey(id), host)
			test.Validator(t, nil, err)
		})
	}
}

func TestGet(t *testing.T) {
	tests := []test.Test{
		{
			Name: "Get should return the host and no error if the session is found",
			Mock: func() {
				redisMock.ExpectGet(MeetingMapKey(id)).SetVal(host)
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Nil(t, err)
				assert.Equal(t, host, value.(string))
			},
		},
		{
			Name: "Get should return an error if redis throws an error",
			Mock: func() {
				mock := redisMock.ExpectGet(MeetingMapKey(id))
				mock.SetErr(errors.New("redis error"))
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			Name: "Get should return an empty string if the session is not found",
			Mock: func() {
				redisMock.ExpectGet(MeetingMapKey(id)).SetVal("")
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
			value, err := mapper.Get(MeetingMapKey(id))
			test.Validator(t, value, err)
		})
	}
}

func TestRemove(t *testing.T) {
	tests := []test.Test{
		{
			Name: "Remove should return nil if the session is removed",
			Mock: func() {
				redisMock.ExpectDel(MeetingMapKey(id)).SetErr(redis.Nil)
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Nil(t, err)
			},
		},
		{
			Name: "Remove should return an error if redis throws an error",
			Mock: func() {
				redisMock.ExpectDel(MeetingMapKey(id)).SetErr(errors.New("redis error"))
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.NotNil(t, err)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Mock()
			err := mapper.Remove(MeetingMapKey(id))
			test.Validator(t, nil, err)
		})
	}
}

func TestDeleteAll(t *testing.T) {
	tests := []test.Test{
		{
			Name: "An error thrown by redis keys method should be returned",
			Mock: func() {
				mock := redisMock.ExpectKeys(RecodingPattern())
				mock.SetErr(errors.New("redis error"))
				mock.SetVal([]string{})
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Error(t, err)
			},
		},
		{
			Name: "An error thrown by redis del method should be returned",
			Mock: func() {
				mock := redisMock.ExpectKeys(RecodingPattern())
				mock.SetVal([]string{RecordingMapKey(id)})
				redisMock.ExpectDel(RecordingMapKey(id)).SetErr(errors.New("redis error"))
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Error(t, err)
			},
		},
		{
			Name: "No error should be returned all keys are deleted",
			Mock: func() {
				mock := redisMock.ExpectKeys(RecodingPattern())
				mock.SetVal([]string{RecordingMapKey(id)})
				redisMock.ExpectDel(RecordingMapKey(id)).SetVal(1)
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Nil(t, err)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Mock()
			err := mapper.DeleteAll(RecodingPattern())
			test.Validator(t, nil, err)
		})
	}
}
