package handler_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"tempo/config"

	"github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	err := config.Load()
	if err != nil {
		fmt.Printf("Config error: %s\n", err.Error())
		os.Exit(1)
	}

	err = initLogging()
	if err != nil {
		fmt.Printf("Logging error: %s\n", err.Error())
		os.Exit(1)
	}

	retCode := m.Run()
	os.Exit(retCode)
}

func initLogging() error {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	log := logrus.StandardLogger()
	level, err := logrus.ParseLevel(config.Instance().LogLevel)
	if err != nil {
		panic(err)
	}
	log.SetLevel(level)

	return err
}

func performRequest(r http.Handler, method string, path string, body io.Reader, headers map[string]string, queryStrings map[string]string) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	if len(queryStrings) > 0 {
		q := req.URL.Query()
		for k, v := range queryStrings {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w, nil
}

func printOnFailed(t *testing.T) func(body string) {
	return func(body string) {
		if t.Failed() {
			t.Logf("Response: %#v", body)
		}
	}
}
