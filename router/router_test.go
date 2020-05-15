package router

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/TipsyPixie/custom-go-webserver/config"
	"github.com/TipsyPixie/custom-go-webserver/httpError"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewFactory(t *testing.T) {
	rootContext := context.Background()

	contextWithoutConfig, cancel := context.WithCancel(rootContext)
	defer cancel()
	_, err := NewFactory(&contextWithoutConfig)
	if err == nil {
		t.Error("context without config succeeded")
	}

	contextWithInvalidTypedConfig := context.WithValue(rootContext, "config", "nothing")
	_, err = NewFactory(&contextWithInvalidTypedConfig)
	if err == nil {
		t.Error("context with invalid type config succeeded")
	}

	testAppConfig := config.Config{
		Env: "test",
		Application: struct {
			Secret string
			Debug  bool
		}{
			Secret: "superDuperSecret",
			Debug:  true,
		},
	}
	contextWithConfig := context.WithValue(rootContext, "config", &testAppConfig)
	_, err = NewFactory(&contextWithConfig)
	if err != nil {
		t.Error(err)
	}
}

func TestHandler_ServeHTTP(t *testing.T) {
	rootContext := context.Background()

	testAppConfig := config.Config{
		Env: "test",
		Application: struct {
			Secret string
			Debug  bool
		}{
			Secret: "superDuperSecret",
			Debug:  true,
		},
	}
	contextWithConfig := context.WithValue(rootContext, "config", &testAppConfig)
	factory, err := NewFactory(&contextWithConfig)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	routeHandler := factory.NewHandler(func(responseWriter http.ResponseWriter, request *http.Request) *httpError.HttpError {
		requestBody, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return httpError.InternalServerError(err, nil)
		}

		if stringBody := string(requestBody); stringBody == "ok" {
			appConfig := request.Context().Value("config").(*config.Config)
			_, err = responseWriter.Write([]byte(appConfig.Application.Secret))
			if err != nil {
				return httpError.InternalServerError(err, nil)
			}
			return nil
		} else {
			return httpError.BadRequest(errors.New(stringBody), nil)
		}
	})

	mockRequest := httptest.NewRequest("GET", "http://test", bytes.NewBufferString("ok"))
	responseRecorder := httptest.NewRecorder()
	routeHandler.ServeHTTP(responseRecorder, mockRequest)
	response := responseRecorder.Result()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if response.StatusCode != http.StatusOK {
		t.Error(fmt.Sprintf("routing handler returns http error code %d", response.StatusCode))
	}
	if secretReturned := string(responseBody); secretReturned != testAppConfig.Application.Secret {
		t.Error(fmt.Sprintf("routing handler returns wrong value %s", secretReturned))
	}

	mockRequest = httptest.NewRequest("GET", "http://test", bytes.NewBufferString("makeError"))
	responseRecorder = httptest.NewRecorder()
	routeHandler.ServeHTTP(responseRecorder, mockRequest)
	response = responseRecorder.Result()
	responseBody, err = ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if response.StatusCode != http.StatusBadRequest {
		t.Error(fmt.Sprintf("routing handler returns http error code %d", response.StatusCode))
	}
	if errorMessage := string(responseBody); errorMessage != "" {
		t.Error(fmt.Sprintf("routing handler returns wrong value %s", errorMessage))
	}
}
