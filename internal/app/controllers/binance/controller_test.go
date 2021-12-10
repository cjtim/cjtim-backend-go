package binance_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cjtim/cjtim-backend-go/configs"

	controllers "github.com/cjtim/cjtim-backend-go/internal/app/controllers"
	"github.com/cjtim/cjtim-backend-go/internal/app/controllers/binance"
	"github.com/cjtim/cjtim-backend-go/internal/app/repository"
	"github.com/cjtim/cjtim-backend-go/internal/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initialBinanceMock() *fiber.App {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", &linebot.UserProfileResponse{
			UserID:      "aaaaaaaaaabbbbbbbbb",
			DisplayName: "Unit Test",
			PictureURL:  "",
		})
		return c.Next()
	})
	app.Get(controllers.RoutePath.Binance.Get, binance.Get)
	app.Get(controllers.RoutePath.Binance.Wallet_GET, binance.GetWallet)
	app.Post(controllers.RoutePath.Binance.Update_POST, binance.UpdatePrice)
	app.Get(controllers.RoutePath.Binance.CronJob_GET, binance.Cronjob)
	return app
}

func Test_Get_FoundUser(t *testing.T) {
	expect := repository.BinanceScheama{
		LineUID: "aaaaaaaaaabbbbbbbbb",
		Prices: map[string]interface{}{
			"BNB": 1.0,
		},
		LineNotifyTime: 5,
	}
	app := initialBinanceMock()

	// Mock
	mock := repository.Mock_Repository{}
	mock.M_FindOne = func(data interface{}, filter interface{}, opts ...*options.FindOneOptions) error {
		b, _ := json.Marshal(&expect)
		return json.Unmarshal(b, data)
	}
	repository.BinanceRepo = &mock
	defer repository.RestoreRepoMock()

	// Test
	req := httptest.NewRequest(http.MethodGet, controllers.RoutePath.Binance.Get, nil)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	defer resp.Body.Close()

	actual := repository.BinanceScheama{}
	err = json.Unmarshal(bodyBytes, &actual)
	assert.Nil(t, err)

	assert.Equal(t, expect, actual)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

}

func Test_Get_NoUser(t *testing.T) {
	// Mock
	expect := repository.BinanceScheama{
		LineUID: "aaaaaaaaaabbbbbbbbb",
		Prices: map[string]interface{}{
			"BNB": 1.00,
		},
		LineNotifyTime: 5,
	}
	mock := repository.Mock_Repository{}
	mock.M_FindOne = func(data interface{}, filter interface{}, opts ...*options.FindOneOptions) error {
		return mongo.ErrNoDocuments
	}
	mock.M_InsertOne = func(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (primitive.ObjectID, error) {
		return primitive.NewObjectID(), nil
	}
	repository.BinanceRepo = &mock
	defer repository.RestoreRepoMock()

	// Test
	app := initialBinanceMock()
	req := httptest.NewRequest(http.MethodGet, controllers.RoutePath.Binance.Get, nil)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	defer resp.Body.Close()
	actual := repository.BinanceScheama{}
	err = json.Unmarshal(bodyBytes, &actual)
	assert.Nil(t, err)

	assert.Equal(t, expect, actual)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

}

func Test_Get_FindOne_Fail(t *testing.T) {
	mock := repository.Mock_Repository{}
	mock.M_FindOne = func(data interface{}, filter interface{}, opts ...*options.FindOneOptions) error {
		return mongo.ErrClientDisconnected
	}
	repository.BinanceRepo = &mock
	defer repository.RestoreRepoMock()

	// Test
	app := initialBinanceMock()
	req := httptest.NewRequest(http.MethodGet, controllers.RoutePath.Binance.Get, nil)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

}

func Test_Get_Insert_Fail(t *testing.T) {
	mock := repository.Mock_Repository{}
	mock.M_FindOne = func(data interface{}, filter interface{}, opts ...*options.FindOneOptions) error {
		return mongo.ErrNoDocuments
	}
	mock.M_InsertOne = func(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (primitive.ObjectID, error) {
		return primitive.NewObjectID(), mongo.ErrClientDisconnected
	}
	repository.BinanceRepo = &mock
	defer repository.RestoreRepoMock()

	// Test
	app := initialBinanceMock()
	req := httptest.NewRequest(http.MethodGet, controllers.RoutePath.Binance.Get, nil)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

}

func Test_GetWallet_Pass(t *testing.T) {
	app := initialBinanceMock()

	mockData := repository.BinanceScheama{
		BinanceApiKey:    "A",
		BinanceSecretKey: "B",
	}
	mock := repository.Mock_Repository{}
	mock.M_FindOne = func(data interface{}, filter interface{}, opts ...*options.FindOneOptions) error {
		b, _ := json.Marshal(&mockData)
		return json.Unmarshal(b, data)
	}
	repository.BinanceRepo = &mock
	defer repository.RestoreRepoMock()

	origBinanceAPI := configs.Config.BinanceAccountAPI
	restoreBinanceAPI := func() { configs.Config.BinanceAccountAPI = origBinanceAPI }
	defer restoreBinanceAPI()

	// start notify microservice server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"balances":[{"asset":"BTC","free":"0.00000000","locked":"0.00000000"}]}`)
	})
	s := httptest.NewServer(handler)
	configs.Config.BinanceAccountAPI = s.URL
	defer s.Close()

	req := httptest.NewRequest(http.MethodGet, controllers.RoutePath.Binance.Wallet_GET, nil)

	resp, err := app.Test(req)
	assert.Nil(t, err)
	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, string(byteBody), `[{"asset":"BTC","free":"0.00000000","locked":"0.00000000"}]`)

}

func Test_GetWallet_BinanceAPI_Fail(t *testing.T) {
	app := initialBinanceMock()

	mockData := repository.BinanceScheama{
		BinanceApiKey:    "A",
		BinanceSecretKey: "B",
	}
	mock := repository.Mock_Repository{}
	mock.M_FindOne = func(data interface{}, filter interface{}, opts ...*options.FindOneOptions) error {
		b, _ := json.Marshal(&mockData)
		return json.Unmarshal(b, data)
	}
	repository.BinanceRepo = &mock
	defer repository.RestoreRepoMock()

	origBinanceAPI := configs.Config.BinanceAccountAPI
	restoreBinanceAPI := func() { configs.Config.BinanceAccountAPI = origBinanceAPI }
	defer restoreBinanceAPI()

	// start notify microservice server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, "")
	})
	s := httptest.NewServer(handler)
	configs.Config.BinanceAccountAPI = s.URL
	defer s.Close()

	req := httptest.NewRequest(http.MethodGet, controllers.RoutePath.Binance.Wallet_GET, nil)

	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

}

func Test_GetWallet_FindOne_Fail(t *testing.T) {
	app := initialBinanceMock()
	mock := repository.Mock_Repository{}
	mock.M_FindOne = func(data interface{}, filter interface{}, opts ...*options.FindOneOptions) error {
		return mongo.ErrNoDocuments
	}
	repository.BinanceRepo = &mock
	defer repository.RestoreRepoMock()

	req := httptest.NewRequest(http.MethodGet, controllers.RoutePath.Binance.Wallet_GET, nil)

	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func Test_GetWallet_NoAPIKey(t *testing.T) {
	app := initialBinanceMock()

	mockData := repository.BinanceScheama{
		BinanceApiKey:    "",
		BinanceSecretKey: "",
	}
	mock := repository.Mock_Repository{}
	mock.M_FindOne = func(data interface{}, filter interface{}, opts ...*options.FindOneOptions) error {
		b, _ := json.Marshal(&mockData)
		return json.Unmarshal(b, data)
	}
	repository.BinanceRepo = &mock
	defer repository.RestoreRepoMock()

	req := httptest.NewRequest(http.MethodGet, controllers.RoutePath.Binance.Wallet_GET, nil)

	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func Test_UpdatePrice_Pass(t *testing.T) {

	expect := repository.BinanceScheama{
		ID:               primitive.ObjectID{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		BinanceSecretKey: "",
		LineNotifyToken:  "",
		BinanceApiKey:    "",
		LineUID:          "aaaaaaaaaabbbbbbbbb",
		Prices: map[string]interface{}{
			"BNB": 1,
		},
		LineNotifyTime: 5,
	}
	mock := repository.Mock_Repository{}
	mock.M_FindOneAndReplace = func(ctx context.Context, filter, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) error {
		return nil
	}
	repository.BinanceRepo = &mock
	defer repository.RestoreRepoMock()

	app := initialBinanceMock()
	req, err := utils.HttpPrepareRequest(&utils.HttpReq{
		Method: http.MethodPost,
		URL:    controllers.RoutePath.Binance.Update_POST,
		Body:   expect,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})
	assert.Nil(t, err)

	resp, err := app.Test(req)
	assert.Nil(t, err)

	defer resp.Body.Close()
	defer assert.Equal(t, http.StatusOK, resp.StatusCode)
}
func Test_UpdatePrice_Fail(t *testing.T) {
	app := initialBinanceMock()

	req, _ := utils.HttpPrepareRequest(&utils.HttpReq{
		Method: http.MethodPost,
		URL:    controllers.RoutePath.Binance.Update_POST,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})

	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func Test_Cron(t *testing.T) {
	mockData := []repository.BinanceScheama{
		{
			BinanceApiKey:    "A",
			BinanceSecretKey: "B",
			LineNotifyToken:  "A",
			LineNotifyTime:   int64(time.Now().Minute()),
		},
		{
			BinanceApiKey:    "AA",
			BinanceSecretKey: "BB",
			LineNotifyToken:  "B",
			LineNotifyTime:   0,
		},
	}
	mock := repository.Mock_Repository{}
	mock.M_Find = func(data, filter interface{}, opts ...*options.FindOptions) error {
		b, _ := json.Marshal(mockData)
		return json.Unmarshal(b, data)
	}
	repository.BinanceRepo = &mock
	defer repository.RestoreRepoMock()

	origLineNotify := configs.Config.LineNotifyURL
	restoreLineNotify := func() { configs.Config.LineNotifyURL = origLineNotify }
	defer restoreLineNotify()

	// start notify microservice server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
	s := httptest.NewServer(handler)
	defer s.Close()
	configs.Config.LineNotifyURL = s.URL

	app := initialBinanceMock()

	req, _ := utils.HttpPrepareRequest(&utils.HttpReq{
		Method: http.MethodGet,
		URL:    controllers.RoutePath.Binance.CronJob_GET,
	})

	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

}

func Test_Cron_Notify_Fail(t *testing.T) {
	mockData := []repository.BinanceScheama{
		{
			BinanceApiKey:    "A",
			BinanceSecretKey: "B",
			LineNotifyToken:  "A",
			LineNotifyTime:   int64(time.Now().Minute()),
		},
		{
			BinanceApiKey:    "AA",
			BinanceSecretKey: "BB",
			LineNotifyToken:  "B",
			LineNotifyTime:   0,
		},
	}
	mock := repository.Mock_Repository{}
	mock.M_Find = func(data, filter interface{}, opts ...*options.FindOptions) error {
		b, _ := json.Marshal(mockData)
		return json.Unmarshal(b, data)
	}
	repository.BinanceRepo = &mock
	defer repository.RestoreRepoMock()

	app := initialBinanceMock()

	req, _ := utils.HttpPrepareRequest(&utils.HttpReq{
		Method: http.MethodGet,
		URL:    controllers.RoutePath.Binance.CronJob_GET,
	})

	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

}

func Test_Cron_Find_Fail(t *testing.T) {
	mock := repository.Mock_Repository{}
	mock.M_Find = func(data, filter interface{}, opts ...*options.FindOptions) error {
		return mongo.ErrNoDocuments
	}
	repository.BinanceRepo = &mock
	defer repository.RestoreRepoMock()

	app := initialBinanceMock()

	req, _ := utils.HttpPrepareRequest(&utils.HttpReq{
		Method: http.MethodGet,
		URL:    controllers.RoutePath.Binance.CronJob_GET,
	})

	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

}
