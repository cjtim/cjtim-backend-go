package binance_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cjtim/cjtim-backend-go/api/binance"
	"github.com/cjtim/cjtim-backend-go/config"
	"github.com/cjtim/cjtim-backend-go/pkg/utils"
	"github.com/cjtim/cjtim-backend-go/repository"
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
	app.Get("/binance/get", binance.Get)
	app.Get("/binance/wallet", binance.GetWallet)
	app.Post("/binance/update", binance.UpdatePrice)
	app.Get("/binance/cronjob", binance.Cronjob)
	return app
}

func Test_Get_FoundUser(t *testing.T) {
	app := initialBinanceMock()
	origFindOne := repository.BinanceRepo.FindOne
	origInsertOne := repository.BinanceRepo.InsertOne

	repository.BinanceRepo.FindOne = func(data interface{}, filter interface{}, opts ...*options.FindOneOptions) error {
		b, _ := json.Marshal(repository.BinanceScheama{
			LineUID: "aaaaaaaaaabbbbbbbbb",
			Prices: map[string]interface{}{
				"BNB": 1,
			},
			LineNotifyTime: 5,
		})
		json.Unmarshal(b, data)
		return nil
	}
	repository.BinanceRepo.InsertOne = func(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (primitive.ObjectID, error) {
		return primitive.NewObjectID(), nil
	}
	req := httptest.NewRequest(http.MethodGet, "/binance/get", nil)

	resp, err := app.Test(req)
	assert.Nil(t, err)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("cannot parse body")
	}
	actual := repository.BinanceScheama{}
	err = json.Unmarshal(bodyBytes, &actual)
	if err != nil {
		t.Fatal(err)
	}
	expect := repository.BinanceScheama{
		LineUID: "aaaaaaaaaabbbbbbbbb",
		Prices: map[string]interface{}{
			"BNB": 1,
		},
		LineNotifyTime: 5,
	}

	assert.Equal(t, expect.LineUID, actual.LineUID)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	repository.BinanceRepo.FindOne = origFindOne
	repository.BinanceRepo.InsertOne = origInsertOne
}

func Test_Get_NoUser(t *testing.T) {
	expect := repository.BinanceScheama{
		LineUID: "aaaaaaaaaabbbbbbbbb",
		Prices: map[string]interface{}{
			"BNB": 1,
		},
		LineNotifyTime: 5,
	}

	app := initialBinanceMock()
	origFindOne := repository.BinanceRepo.FindOne
	origInsertOne := repository.BinanceRepo.InsertOne

	findOneCount := 0
	repository.BinanceRepo.FindOne = func(data interface{}, filter interface{}, opts ...*options.FindOneOptions) error {
		if findOneCount == 0 {
			findOneCount++
			return mongo.ErrNoDocuments
		}
		b, _ := json.Marshal(expect)
		return json.Unmarshal(b, data)
	}
	repository.BinanceRepo.InsertOne = func(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (primitive.ObjectID, error) {
		return primitive.NewObjectID(), nil
	}
	req := httptest.NewRequest(http.MethodGet, "/binance/get", nil)

	resp, err := app.Test(req)
	assert.Nil(t, err)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("cannot parse body")
	}
	actual := repository.BinanceScheama{}
	err = json.Unmarshal(bodyBytes, &actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expect.LineUID, actual.LineUID)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	repository.BinanceRepo.FindOne = origFindOne
	repository.BinanceRepo.InsertOne = origInsertOne
}

func Test_GetWallet_Pass(t *testing.T) {
	app := initialBinanceMock()

	origFindOne := repository.BinanceRepo.FindOne

	repository.BinanceRepo.FindOne = func(data interface{}, filter interface{}, opts ...*options.FindOneOptions) error {
		return nil
	}

	req := httptest.NewRequest(http.MethodGet, "/binance/wallet", nil)

	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	repository.BinanceRepo.FindOne = origFindOne
}

func Test_UpdatePrice_Pass(t *testing.T) {
	orgi := repository.BinanceRepo.FindOneAndReplace
	origFindOneAndReplace := func() {
		repository.BinanceRepo.FindOneAndReplace = orgi
	}

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
	repository.BinanceRepo.FindOneAndReplace = func(ctx context.Context, filter, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) error {
		return nil
	}

	app := initialBinanceMock()
	req, err := utils.HttpPrepareRequest(&utils.HttpReq{
		Method: http.MethodPost,
		URL:    "/binance/update",
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
	defer origFindOneAndReplace()
}
func Test_UpdatePrice_Fail(t *testing.T) {
	app := initialBinanceMock()

	req, _ := utils.HttpPrepareRequest(&utils.HttpReq{
		Method: http.MethodPost,
		URL:    "/binance/update",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})

	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func Test_Cron(t *testing.T) {
	orig := repository.BinanceRepo.Find
	origFind := func() {
		repository.BinanceRepo.Find = orig
	}
	repository.BinanceRepo.Find = func(data, filter interface{}, opts ...*options.FindOptions) error {
		return nil
	}

	origLineNotify := config.Config.LineNotifyURL
	restoreLineNotify := func() { config.Config.LineNotifyURL = origLineNotify }

	// start notify microservice server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, "")
	})
	s := httptest.NewServer(handler)
	defer s.Close()
	config.Config.LineNotifyURL = s.URL

	app := initialBinanceMock()

	req, _ := utils.HttpPrepareRequest(&utils.HttpReq{
		Method: http.MethodGet,
		URL:    "/binance/cronjob",
	})

	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	defer restoreLineNotify()
	defer origFind()
}
