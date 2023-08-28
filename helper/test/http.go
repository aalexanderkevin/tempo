package test

import (
	"tempo/config"
	"tempo/container"
	"tempo/controller"

	"net/http"
	"testing"
)

func SetupHttpHandler(t *testing.T, cb func(appContainer *container.Container) *container.Container) http.Handler {
	appContainer := DefaultAppContainer()
	if cb != nil {
		appContainer = cb(appContainer)
	}
	server := controller.NewHttpServer(appContainer)
	handler, err := server.GetHandler()
	if err != nil {
		t.Fatal(err.Error())
	}
	return handler
}

func DefaultAppContainer() *container.Container {
	appContainer := container.NewContainer()

	appContainer.SetConfig(config.Instance())

	return appContainer
}
