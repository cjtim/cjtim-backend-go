package main

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2/utils"
)

func init() {
}

var app = startServer()

func Test_Route_Home(t *testing.T) {
	resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
	utils.AssertEqual(t, nil, err, "is error?")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
	body, err := ioutil.ReadAll(resp.Body)
	utils.AssertEqual(t, "{\"msg\":\"Hello, world\"}", string(body), "hello world")

}

func Test_Route_Ping(t *testing.T) {
	resp, err := app.Test(httptest.NewRequest("GET", "/ping", nil))
	utils.AssertEqual(t, nil, err, "is error?")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
	body, err := ioutil.ReadAll(resp.Body)
	utils.AssertEqual(t, "pong", string(body), "PING PONG")
}

func Test_Route_Me(t *testing.T) {
	resp, err := app.Test(httptest.NewRequest("GET", "/users/me", nil))
	utils.AssertEqual(t, nil, err, "is error?")
	utils.AssertEqual(t, 403, resp.StatusCode, "Status code: 403")
	body, err := ioutil.ReadAll(resp.Body)
	utils.AssertEqual(t, "Forbidden", string(body), "Forbidden access due to no Authorization headers")
}

func Test_Shutdown_Server(t *testing.T) {
	app.Server().Shutdown()
}
