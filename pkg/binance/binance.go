package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/cjtim/cjtim-backend-go/config"
	"github.com/cjtim/cjtim-backend-go/pkg/utils"
	"go.uber.org/zap"
)

func ComputeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func GetBinanceAccount(apiKey string, secretKey string) (map[string]interface{}, error) {
	timeNow := time.Now().UnixNano() / int64(time.Millisecond)
	signature := ComputeHmac256("timestamp="+fmt.Sprint(timeNow), secretKey)
	resp, respBody, err := utils.Http(&utils.HttpReq{
		Method:  http.MethodGet,
		URL:     config.Config.BinanceAccountAPI,
		Headers: map[string]string{"X-MBX-APIKEY": apiKey},
		Querys: map[string]string{
			"timestamp": fmt.Sprint(timeNow),
			"signature": signature,
		},
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		err := errors.New(string(respBody))
		zap.L().Info("GetBinanceAccount", zap.Error(err))
		return nil, err
	}

	binanceAccount := map[string]interface{}{}
	err = json.Unmarshal(respBody, &binanceAccount)
	if err != nil {
		zap.L().Info("GetBinanceAccount", zap.Error(err))
		return nil, err
	}
	return binanceAccount, nil
}
