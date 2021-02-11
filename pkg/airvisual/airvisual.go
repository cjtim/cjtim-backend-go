package airvisual

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

var _ = godotenv.Load()
var restyClient = resty.New()

func get(queryParams map[string]string, targetAPI string) (*AirVisualResponse, error) {
	resp, err := restyClient.R().SetQueryParams(queryParams).Get(targetAPI)
	body := &AirVisualResponse{}
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, errors.New(string(resp.Body()))
	}
	err = json.Unmarshal(resp.Body(), body)
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
