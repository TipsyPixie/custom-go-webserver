package router

import (
	"context"
	"custom-go-webserver/config"
	"custom-go-webserver/httpError"
	"errors"
	"net/http"
)

type ContextKey string

type HandleFunc func(responseWriter http.ResponseWriter, request *http.Request) *httpError.HttpError

type Handler struct {
	handle     HandleFunc
	appContext *context.Context
}

type HandlerFactory struct {
	appContext *context.Context
}

func NewFactory(appContext *context.Context) (*HandlerFactory, error) {
	appConfig := (*appContext).Value("config")
	switch appConfig.(type) {
	case *config.Config:
	default:
		return nil, errors.New("config type must be *config.Config")
	}

	return &HandlerFactory{
		appContext: appContext,
	}, nil
}

func (thisHandlerFactory *HandlerFactory) NewHandler(f HandleFunc) *Handler {
	return &Handler{
		handle:     f,
		appContext: thisHandlerFactory.appContext,
	}
}

func (thisHandler *Handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	requestContext, cancelRequestContext := context.WithCancel(*thisHandler.appContext)
	defer cancelRequestContext()
	requestWithContext := request.Clone(requestContext)
	httpErr := thisHandler.handle(responseWriter, requestWithContext)
	if httpErr != nil {
		responseWriter.WriteHeader(httpErr.StatusCode)
		_, _ = responseWriter.Write(httpErr.Description)
	}
}
