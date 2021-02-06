package main

import (
	"context"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/cjtim/cjtim-backend-go/datasource/collections"
	"github.com/cjtim/cjtim-backend-go/models"
	"github.com/gofiber/fiber/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
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

func Test_Database(t *testing.T) {
	m, err := models.GetModels(nil)
	if err != nil {
		t.Fatal(err)
	}

	newData := bson.M{
		"test": "123567",
	}
	id, err := m.InsertOne("unit_test", newData)
	if err != nil {
		t.Fatal("Failed Insert DB!")
	}
	actual := &collections.UnitTestSchema{}
	m.FindOne("unit_test", &actual, newData)
	if actual.Test != newData["test"] {
		t.Fatal("Failed FindOne DB!")
	}
	if m.Destroy("unit_test", bson.M{"_id": id}) != nil {
		t.Fatal("Failed Destroy DB!")
	}
	if m.Client.Disconnect(context.TODO()) != nil {
		t.Fatal("Failed Disconect DB!")
	}
}

func Test_Shutdown_Server(t *testing.T) {
	app.Server().Shutdown()
}
