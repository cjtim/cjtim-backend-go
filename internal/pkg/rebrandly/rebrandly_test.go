package rebrandly_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cjtim/cjtim-backend-go/configs"
	"github.com/cjtim/cjtim-backend-go/internal/app/repository"
	"github.com/cjtim/cjtim-backend-go/internal/pkg/rebrandly"
	"github.com/stretchr/testify/assert"
)

func Test_Add_Success(t *testing.T) {
	expect := &repository.URLScheama{
		RebrandlyID: "a",
		ShortURL:    "link.cjtim.com/abc",
		Destination: "google.com",
	}
	bExpect, _ := json.Marshal(expect)
	// start notify microservice server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(bExpect))
		w.WriteHeader(http.StatusOK)
	})
	s := httptest.NewServer(handler)
	defer s.Close()
	configs.Config.RebrandlyAPI = s.URL
	defer configs.RestoreConfigMock()

	actual, err := rebrandly.Add("google.com")
	assert.Nil(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expect, actual)
}

func Test_Add_API_Fail(t *testing.T) {
	// start notify microservice server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
	})
	s := httptest.NewServer(handler)
	defer s.Close()
	configs.Config.RebrandlyAPI = s.URL
	defer configs.RestoreConfigMock()

	actual, err := rebrandly.Add("google.com")
	assert.NotNil(t, err)
	assert.Nil(t, actual)
}

func Test_Add_API_Invalid_URL(t *testing.T) {
	configs.Config.RebrandlyAPI = "Fake api url"
	defer configs.RestoreConfigMock()

	actual, err := rebrandly.Add("google.com")
	assert.NotNil(t, err)
	assert.Nil(t, actual)
}

func Test_Delete_Success(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
	s := httptest.NewServer(handler)
	defer s.Close()
	configs.Config.RebrandlyAPI = s.URL
	defer configs.RestoreConfigMock()

	err := rebrandly.Delete("1")
	assert.Nil(t, err)
}

func Test_Delete_BadRequest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
	})
	s := httptest.NewServer(handler)
	defer s.Close()
	configs.Config.RebrandlyAPI = s.URL
	defer configs.RestoreConfigMock()

	err := rebrandly.Delete("1")
	assert.NotNil(t, err)
}
func Test_Delete_Invalid_URL(t *testing.T) {
	configs.Config.RebrandlyAPI = "Fake api"
	defer configs.RestoreConfigMock()

	err := rebrandly.Delete("1")
	assert.NotNil(t, err)
}
