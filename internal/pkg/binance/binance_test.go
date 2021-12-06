package binance_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cjtim/cjtim-backend-go/configs"
	"github.com/cjtim/cjtim-backend-go/internal/pkg/binance"
	"github.com/stretchr/testify/assert"
)

func Test_GetBinanceAccount_Success(t *testing.T) {
	origBinanceAPI := configs.Config.BinanceAccountAPI
	restoreBinanceAPI := func() { configs.Config.BinanceAccountAPI = origBinanceAPI }
	defer restoreBinanceAPI()

	// start notify microservice server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, "{\"test\":\"test\"}")
	})
	s := httptest.NewServer(handler)
	configs.Config.BinanceAccountAPI = s.URL
	defer s.Close()

	result, err := binance.GetBinanceAccount("a", "b")
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{
		"test": "test",
	}, result)
}

func Test_GetBinanceAccount_API_Fail(t *testing.T) {
	origBinanceAPI := configs.Config.BinanceAccountAPI
	restoreBinanceAPI := func() { configs.Config.BinanceAccountAPI = origBinanceAPI }
	defer restoreBinanceAPI()

	// start notify microservice server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// fmt.Fprintln(w, "{\"test\":\"test\"}")
		w.WriteHeader(http.StatusInternalServerError)
	})
	s := httptest.NewServer(handler)
	configs.Config.BinanceAccountAPI = s.URL
	defer s.Close()

	result, err := binance.GetBinanceAccount("a", "b")
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func Test_GetBinanceAccount_API_InvalidBody(t *testing.T) {
	origBinanceAPI := configs.Config.BinanceAccountAPI
	restoreBinanceAPI := func() { configs.Config.BinanceAccountAPI = origBinanceAPI }
	defer restoreBinanceAPI()

	// start notify microservice server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, "{\"test\":\"te")
		w.WriteHeader(http.StatusOK)
	})
	s := httptest.NewServer(handler)
	configs.Config.BinanceAccountAPI = s.URL
	defer s.Close()

	result, err := binance.GetBinanceAccount("a", "b")
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func Test_GetBinanceAccount_API_InvalidEndpoint(t *testing.T) {
	origBinanceAPI := configs.Config.BinanceAccountAPI
	restoreBinanceAPI := func() { configs.Config.BinanceAccountAPI = origBinanceAPI }
	defer restoreBinanceAPI()

	// start notify microservice server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// fmt.Fprintln(w, "{\"test\":\"test\"}")
		w.WriteHeader(http.StatusInternalServerError)
	})
	s := httptest.NewServer(handler)
	configs.Config.BinanceAccountAPI = "Fake endpoint"
	defer s.Close()

	result, err := binance.GetBinanceAccount("a", "b")
	assert.NotNil(t, err)
	assert.Nil(t, result)
}
