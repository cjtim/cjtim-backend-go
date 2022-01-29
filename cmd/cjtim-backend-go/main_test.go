package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cjtim/cjtim-backend-go/internal/app/repository"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/stretchr/testify/assert"
)

func Test_Main_Fail_DB(t *testing.T) {
	a := realMain()
	assert.Equal(t, 1, a)
}

func Test_Route_Home(t *testing.T) {
	app := startServer()
	resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
	utils.AssertEqual(t, nil, err, "is error?")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
	body, err := ioutil.ReadAll(resp.Body)
	utils.AssertEqual(t, "{\"msg\":\"Hello, world\"}", string(body), "hello world")
	utils.AssertEqual(t, nil, err)
}

func Test_Route_Ping(t *testing.T) {
	app := startServer()

	// Mock
	mock := repository.Mock_Repository{}
	mock.M_Health = func() error {
		return nil
	}
	repository.BinanceRepo = &mock
	repository.FileRepo = &mock
	repository.UrlRepo = &mock
	repository.UserRepo = &mock
	defer repository.RestoreRepoMock()

	resp, err := app.Test(httptest.NewRequest("GET", "/health", nil))
	utils.AssertEqual(t, nil, err, "is error?")
	utils.AssertEqual(t, http.StatusOK, resp.StatusCode)
}

func Test_Route_Me(t *testing.T) {
	app := startServer()

	resp, err := app.Test(httptest.NewRequest("GET", "/users/me", nil))
	utils.AssertEqual(t, nil, err, "is error?")
	utils.AssertEqual(t, 403, resp.StatusCode, "Status code: 403")
	body, err := ioutil.ReadAll(resp.Body)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, "Forbidden", string(body), "Forbidden access due to no Authorization headers")
}

func Test_Shutdown_Server(t *testing.T) {
	app := startServer()
	app.Server().Shutdown()
}

func Test_CloseHandler(t *testing.T) {
	app := startServer()
	setupCloseHandler(app)
}
