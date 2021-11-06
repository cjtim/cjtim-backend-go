package airvisual

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/cjtim/cjtim-backend-go/pkg/utils"
	"github.com/joho/godotenv"
)

var _ = godotenv.Load()

func get(queryParams map[string]string, targetAPI string) (*AirVisualResponse, error) {
	resp, respBody, err := utils.HttpGET(targetAPI, queryParams, nil)
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
		"key": os.Getenv("AIR_API_KEY"),
	}, "http://api.airvisual.com/v2/nearest_city")
}
func GetPhayaThaiCity() (*AirVisualResponse, error) {
	return get(map[string]string{
		"city":    "Phaya Thai",
		"state":   "Bangkok",
		"country": "Thailand",
		"key":     os.Getenv("AIR_API_KEY"),
	}, "http://api.airvisual.com/v2/city")
}
