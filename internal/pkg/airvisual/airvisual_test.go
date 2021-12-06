package airvisual_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cjtim/cjtim-backend-go/configs"
	"github.com/cjtim/cjtim-backend-go/internal/pkg/airvisual"
	"github.com/stretchr/testify/assert"
)

func Test_ByLocation_Pass(t *testing.T) {
	expect, _ := json.Marshal(airvisual.AirVisualResponse{
		Status: "success",
	})
	origAPI := configs.Config.AirVisualAPINearestCity
	restoreAPI := func() {
		configs.Config.AirVisualAPINearestCity = origAPI
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "{\"status\":\"success\"}")
	})
	s := httptest.NewServer(handler)
	defer s.Close()
	configs.Config.AirVisualAPINearestCity = s.URL
	defer restoreAPI()

	resp, err := airvisual.GetByLocation(1.0, 2.0)
	assert.Nil(t, err)
	assert.Equal(t, "success", resp.Status)

	actual, _ := json.Marshal(resp)
	assert.Equal(t, expect, actual)
}

func Test_ByLocation_InvalidBody(t *testing.T) {
	origAPI := configs.Config.AirVisualAPINearestCity
	restoreAPI := func() {
		configs.Config.AirVisualAPINearestCity = origAPI
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "{\"status\":\"success\"")
	})
	s := httptest.NewServer(handler)
	defer s.Close()
	configs.Config.AirVisualAPINearestCity = s.URL
	defer restoreAPI()

	resp, err := airvisual.GetByLocation(1.0, 2.0)
	assert.NotNil(t, err)
	assert.Nil(t, nil, resp)

}

func Test_ByLocation_Fail(t *testing.T) {
	origAPI := configs.Config.AirVisualAPINearestCity
	restoreAPI := func() {
		configs.Config.AirVisualAPINearestCity = origAPI
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
	})
	s := httptest.NewServer(handler)
	defer s.Close()
	configs.Config.AirVisualAPINearestCity = s.URL
	defer restoreAPI()

	resp, err := airvisual.GetByLocation(1.0, 2.0)
	assert.Equal(t, "", err.Error())
	assert.Nil(t, nil, resp)
}

func Test_City_Pass(t *testing.T) {
	expect, _ := json.Marshal(airvisual.AirVisualResponse{
		Status: "success",
	})
	origAPI := configs.Config.AirVisualAPICity
	restoreAPI := func() {
		configs.Config.AirVisualAPICity = origAPI
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "{\"status\":\"success\"}")
	})
	s := httptest.NewServer(handler)
	defer s.Close()
	configs.Config.AirVisualAPICity = s.URL
	defer restoreAPI()

	resp, err := airvisual.GetPhayaThaiCity()
	assert.Nil(t, err)
	assert.Equal(t, "success", resp.Status)

	actual, _ := json.Marshal(resp)
	assert.Equal(t, expect, actual)
}

func Test_City_InvalidBody(t *testing.T) {
	origAPI := configs.Config.AirVisualAPICity
	restoreAPI := func() {
		configs.Config.AirVisualAPICity = origAPI
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "{\"status\":\"success\"")
	})
	s := httptest.NewServer(handler)
	defer s.Close()
	configs.Config.AirVisualAPICity = s.URL
	defer restoreAPI()

	resp, err := airvisual.GetPhayaThaiCity()
	assert.NotNil(t, err)
	assert.Nil(t, nil, resp)

}

func Test_City_Fail(t *testing.T) {
	origAPI := configs.Config.AirVisualAPICity
	restoreAPI := func() {
		configs.Config.AirVisualAPICity = origAPI
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
	})
	s := httptest.NewServer(handler)
	defer s.Close()
	configs.Config.AirVisualAPICity = s.URL
	defer restoreAPI()

	resp, err := airvisual.GetPhayaThaiCity()
	assert.Equal(t, "", err.Error())
	assert.Nil(t, nil, resp)
}

func Test_API_InvalidEndpoint(t *testing.T) {
	origAPI := configs.Config.AirVisualAPICity
	restoreAPI := func() {
		configs.Config.AirVisualAPICity = origAPI
	}

	configs.Config.AirVisualAPICity = "Fake endpoint"
	defer restoreAPI()

	resp, err := airvisual.GetPhayaThaiCity()
	assert.NotNil(t, err)
	assert.Nil(t, nil, resp)
}
