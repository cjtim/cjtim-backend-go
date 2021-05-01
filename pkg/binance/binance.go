package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

var restyClient = resty.New()

func ComputeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func GetBinanceAccount(apiKey string, secretKey string) (map[string]interface{}, error) {
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
