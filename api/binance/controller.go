package binance

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/cjtim/cjtim-backend-go/datasource/collections"
	"github.com/cjtim/cjtim-backend-go/models"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/line/line-bot-sdk-go/linebot"
	"go.mongodb.org/mongo-driver/bson"
)

var restyClient = resty.New()

func Get(c *fiber.Ctx) error {
	data := collections.BinanceScheama{}
	user := c.Locals("user").(*linebot.UserProfileResponse)
	models := c.Locals("db").(*models.Models)
	collection := models.Client.Database("production").Collection("binance")
	result := collection.FindOne(context.TODO(), bson.M{"lineUid": user.UserID})
	err := result.Decode(&data)
	if err != nil {
		newUser := collections.BinanceScheama{
			LineUID:          user.UserID,
			LineNotifyToken:  "",
			BinanceApiKey:    "",
			BinanceSecretKey: "",
			Prices:           map[string]interface{}{},
		}
		collection.InsertOne(context.TODO(), newUser)
		result := collection.FindOne(context.TODO(), bson.M{"lineUid": user.UserID})
		err := result.Decode(&data)
		if err != nil {
			return nil
		}
	}
	respData := map[string]interface{}{}
	dataByte, _ := json.Marshal(data)
	json.Unmarshal(dataByte, &respData)
	if data.BinanceApiKey != "" && data.BinanceSecretKey != "" {
		binanceAccount, err := getBinanceAccount(data.BinanceApiKey, data.BinanceSecretKey)
		if err != nil {
			return err
		}
		respData["balances"] = binanceAccount["balances"]
	}
	return c.JSON(respData)
}
func UpdatePrice(c *fiber.Ctx) error {
	data := collections.BinanceScheama{}
	err := c.BodyParser(&data)
	if err != nil {
		return err
	}
	user := c.Locals("user").(*linebot.UserProfileResponse)
	models := c.Locals("db").(*models.Models)
	collection := models.Client.Database("production").Collection("binance")
	collection.FindOneAndReplace(context.TODO(), bson.M{"lineUid": user.UserID}, data)
	return c.SendStatus(200)
}

func ComputeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func getBinanceAccount(apiKey string, secretKey string) (map[string]interface{}, error) {
	timeNow := time.Now().UnixNano() / int64(time.Millisecond)
	signature := ComputeHmac256("timestamp="+fmt.Sprint(timeNow), secretKey)
	url := "https://api.binance.com/api/v3/account?timestamp=" + fmt.Sprint(timeNow)
	url += "&signature=" + signature
	resp, err := restyClient.R().SetHeader("X-MBX-APIKEY", apiKey).Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, errors.New(string(resp.Body()))
	}
	binanceAccount := map[string]interface{}{}
	err = json.Unmarshal(resp.Body(), &binanceAccount)
	if err != nil {
		return nil, err
	}
	return binanceAccount, nil
}
