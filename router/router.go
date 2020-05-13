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
	handle        HandleFunc
	appContext    *context.Context
	preProcessors []HandleFunc
}

func (thisHandler *Handler) BeforeHandler(f HandleFunc) *Handler {
	thisHandler.preProcessors = append(thisHandler.preProcessors, f)
	return thisHandler
}

type HandlerFactory struct {
	appContext    *context.Context
	preProcessors []HandleFunc
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

func (thisHandlerFactory *HandlerFactory) BeforeHandler(f HandleFunc) *HandlerFactory {
	thisHandlerFactory.preProcessors = append(thisHandlerFactory.preProcessors, f)
	return thisHandlerFactory
}

func (thisHandler *Handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	requestContext, cancelRequestContext := context.WithCancel(*thisHandler.appContext)
	defer cancelRequestContext()
	*request = *(request.Clone(requestContext))

	for _, preProcessor := range thisHandler.preProcessors {
		httpErr := preProcessor(responseWriter, request)
		if httpErr != nil {
			responseWriter.WriteHeader(httpErr.StatusCode)
			_, _ = responseWriter.Write(httpErr.Description)
			return
		}
	}

	httpErr := thisHandler.handle(responseWriter, request)
	if httpErr != nil {
		responseWriter.WriteHeader(httpErr.StatusCode)
		_, _ = responseWriter.Write(httpErr.Description)
		return
	}
}
