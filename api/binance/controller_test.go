package binance_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cjtim/cjtim-backend-go/api/binance"
	"github.com/cjtim/cjtim-backend-go/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	client, err := repository.MongoClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	repository.DB = client.Database("production")
	repository.Client = client
}

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

func Test_Get_NoUser(t *testing.T) {
	app := initialBinanceMock()

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
}

func Test_GetWallet_Pass(t *testing.T) {
	app := initialBinanceMock()

	req := httptest.NewRequest(http.MethodGet, "/binance/wallet", nil)

	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_UpdatePrice_Pass(t *testing.T) {
	data := repository.BinanceScheama{}
	result := repository.DB.Collection("binance").FindOne(context.TODO(), bson.M{"lineUid": "aaaaaaaaaabbbbbbbbb"})
	err := result.Decode(&data)
	assert.Nil(t, err)
	byteBody, err := json.Marshal(data)
	assert.Nil(t, err)

	app := initialBinanceMock()

	req := httptest.NewRequest(http.MethodPost, "/binance/update", strings.NewReader(string(byteBody)))

	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}
func Test_UpdatePrice_Fail(t *testing.T) {
	app := initialBinanceMock()

	req := httptest.NewRequest(http.MethodPost, "/binance/update", nil)

	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}

func Test_Cron(t *testing.T) {
	needNotify := false
	user := repository.BinanceScheama{
		LineNotifyTime: int64(time.Now().Minute()),
	}
	userTime := (user.LineNotifyTime) % 60
	currentMinute := time.Now().Minute()
	if userTime > 0 {
		needNotify = (currentMinute % int(userTime)) == 0
	} else {
		needNotify = userTime == int64(currentMinute)
	}
	t.Log(currentMinute)
	if needNotify {
		t.Log(currentMinute)
	} else {
		t.Fatal("Not right time")
	}
}
