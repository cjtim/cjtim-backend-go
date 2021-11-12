package airvisual

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cjtim/cjtim-backend-go/config"
	"github.com/cjtim/cjtim-backend-go/pkg/utils"
)

func get(queryParams map[string]string, targetAPI string) (*AirVisualResponse, error) {
	resp, respBody, err := utils.Http(&utils.HttpReq{
		Method: http.MethodGet,
		URL:    targetAPI,
		Querys: queryParams,
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(respBody))
	}

	body := &AirVisualResponse{}
	err = json.Unmarshal(respBody, body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetByLocation(lat float64, lon float64) (*AirVisualResponse, error) {
	return get(map[string]string{
		"lat": fmt.Sprintf("%f", lat),
		"lon": fmt.Sprintf("%f", lon),
		"key": config.Config.AirVisualKey,
	}, config.Config.AirVisualAPINearestCity)
}
func GetPhayaThaiCity() (*AirVisualResponse, error) {
	return get(map[string]string{
		"city":    "Phaya Thai",
		"state":   "Bangkok",
		"country": "Thailand",
		"key":     config.Config.AirVisualKey,
	}, config.Config.AirVisualAPICity)
}
