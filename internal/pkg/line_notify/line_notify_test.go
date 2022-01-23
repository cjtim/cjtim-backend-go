package line_notify_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cjtim/cjtim-backend-go/configs"
	"github.com/cjtim/cjtim-backend-go/internal/app/repository"
	"github.com/cjtim/cjtim-backend-go/internal/pkg/line_notify"
	"github.com/stretchr/testify/assert"
)

func Test_Success(t *testing.T) {
	mockData := []repository.BinanceScheama{
		{
			BinanceApiKey:    "A",
			BinanceSecretKey: "B",
			LineNotifyToken:  "A",
			LineNotifyTime:   int64(time.Now().Minute()),
			LineUID:          "A",
		},
		{
			BinanceApiKey:    "AA",
			BinanceSecretKey: "BB",
			LineNotifyToken:  "B",
			LineNotifyTime:   int64(time.Now().Minute()),
			LineUID:          "B",
		},
	}
	// start notify microservice server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
	s := httptest.NewServer(handler)
	defer s.Close()
	configs.Config.LineNotifyURL = s.URL

	sUser, eUser := line_notify.TriggerLineNotify(&mockData)
	assert.ElementsMatch(t, []string{"A", "B"}, sUser)
	assert.Equal(t, []string{}, eUser)
}

func Test_FailSomeUser(t *testing.T) {
	mockData := []repository.BinanceScheama{
		{
			BinanceApiKey:    "A",
			BinanceSecretKey: "B",
			LineNotifyToken:  "A",
			LineNotifyTime:   int64(time.Now().Minute()),
			LineUID:          "A",
		},
		{
			BinanceApiKey:    "AA",
			BinanceSecretKey: "BB",
			LineNotifyToken:  "B",
			LineNotifyTime:   int64(time.Now().Minute()),
			LineUID:          "B",
		},
	}

	requestCount := 0
	// start notify microservice server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		requestCount++
		if requestCount > 1 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	s := httptest.NewServer(handler)
	defer s.Close()
	configs.Config.LineNotifyURL = s.URL

	sUser, eUser := line_notify.TriggerLineNotify(&mockData)
	assert.Len(t, sUser, 1)
	assert.Len(t, eUser, 1)
	assert.NotEqual(t, sUser, eUser)
}
