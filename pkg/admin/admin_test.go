package admin

import (
	"b3lb/internal/test"
	"b3lb/pkg/admin/mock"
	"b3lb/pkg/api"
	"b3lb/pkg/config"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAddInstance(t *testing.T) {
	var w *httptest.ResponseRecorder
	var c *gin.Context
	admin := CreateAdmin(&mock.InstanceManager{}, &config.AdminConfig{})
	tests := []test.Test{
		{
			Name: "Add should return a bad request if the request body is empty",
			Mock: func() {
				test.AddRequestBody(c, "")
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			Name: "Add should return an internal server error if the instance manager returns an error when checking if instance already exists",
			Mock: func() {
				mock.ExistsFunc = func(instance api.BigBlueButtonInstance) (bool, error) {
					return false, errors.New("unexpected error")
				}
				test.SetRequestMethod(c, http.MethodPost)
				test.SetRequestContentType(c, "application/json")
				test.AddRequestBody(c, `{"url":"http://localhost/bigbluebutton", "secret":"secret"}`)
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			Name: "Add should return a conflict error if instance already exists",
			Mock: func() {
				mock.ExistsFunc = func(instance api.BigBlueButtonInstance) (bool, error) {
					return true, nil
				}
				test.SetRequestMethod(c, http.MethodPost)
				test.SetRequestContentType(c, "application/json")
				test.AddRequestBody(c, `{"url":"http://localhost/bigbluebutton", "secret":"secret"}`)
			},
			Validator: func(vt *testing.T, alue interface{}, err error) {
				assert.Equal(t, http.StatusConflict, w.Code)
			},
		},
		{
			Name: "Add should return an internal server error if instance manager returns an error when adding instance",
			Mock: func() {
				mock.ExistsFunc = func(instance api.BigBlueButtonInstance) (bool, error) {
					return false, nil
				}
				mock.AddFunc = func(instance api.BigBlueButtonInstance) error {
					return errors.New("unexpected error")
				}
				test.SetRequestMethod(c, http.MethodPost)
				test.SetRequestContentType(c, "application/json")
				test.AddRequestBody(c, `{"url":"http://localhost/bigbluebutton", "secret":"secret"}`)
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			Name: "Add should return a created status and instance if instance is added successfully",
			Mock: func() {
				mock.ExistsFunc = func(instance api.BigBlueButtonInstance) (bool, error) {
					return false, nil
				}
				mock.AddFunc = func(instance api.BigBlueButtonInstance) error {
					return nil
				}
				test.SetRequestMethod(c, http.MethodPost)
				test.SetRequestContentType(c, "application/json")
				test.AddRequestBody(c, `{"url":"http://localhost/bigbluebutton", "secret":"secret"}`)
			},
			Validator: func(vt *testing.T, alue interface{}, err error) {
				assert.Equal(t, w.Body.String(), `{"url":"http://localhost/bigbluebutton","secret":"secret"}`)
				assert.Equal(t, http.StatusCreated, w.Code)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)
			test.Mock()
			admin.AddInstance(c)
			test.Validator(t, nil, nil)
		})
	}
}

func TestListInstances(t *testing.T) {
	url := "http://localhost/bigbluebutton"
	var w *httptest.ResponseRecorder
	admin := CreateAdmin(&mock.InstanceManager{}, &config.AdminConfig{})

	tests := []test.Test{
		{
			Name: "List should returns a list containg a single bigbluebutton instance",
			Mock: func() {
				mock.ListFunc = func() ([]string, error) {
					return []string{url}, nil
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, test.StringToJSONArray(w.Body.String())[0], url)
			},
		},
		{
			Name: "List should return an internal server error if instance manager returns an error",
			Mock: func() {
				mock.ListFunc = func() ([]string, error) {
					return nil, errors.New("unexpected error")
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			test.Mock()
			admin.ListInstances(c)
			test.Validator(t, nil, nil)
		})
	}
}

func TestDeleteInstance(t *testing.T) {
	var w *httptest.ResponseRecorder
	var c *gin.Context
	admin := CreateAdmin(&mock.InstanceManager{}, &config.AdminConfig{})
	tests := []test.Test{
		{
			Name: "Delete should return a bad request status if instance url is not provided",
			Mock: func() {},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			Name: "Delete should return an internal server error if instance manager returns an error",
			Mock: func() {
				mock.ExistsFunc = func(instance api.BigBlueButtonInstance) (bool, error) {
					return false, errors.New("unexpected error")
				}
				test.SetRequestParams(c, "url=http://localhost/bigbluebutton")
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			Name: "Delete should return a not found error if instance does not exist",
			Mock: func() {
				mock.ExistsFunc = func(instance api.BigBlueButtonInstance) (bool, error) {
					return false, nil
				}
				test.SetRequestParams(c, "url=http://localhost/bigbluebutton")
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		{
			Name: "Delete should an internal server error if instance manager returns an error on deleting instance",
			Mock: func() {
				mock.ExistsFunc = func(instance api.BigBlueButtonInstance) (bool, error) {
					return true, nil
				}
				mock.RemoveFunc = func(URL string) error {
					return errors.New("unexpected error")
				}
				test.SetRequestParams(c, "url=http://localhost/bigbluebutton")
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			Name: "Delete should return a no content status if instance is deleted",
			Mock: func() {
				mock.ExistsFunc = func(instance api.BigBlueButtonInstance) (bool, error) {
					return true, nil
				}
				mock.RemoveFunc = func(URL string) error {
					return nil
				}
				test.SetRequestParams(c, "url=http://localhost/bigbluebutton")
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusNoContent, w.Code)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)
			test.Mock()
			admin.DeleteInstance(c)
			test.Validator(t, nil, nil)
		})
	}
}
