package airvisual

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

var _ = godotenv.Load()
var restyClient = resty.New()

func Get(queryParams map[string]string) (*AirVisualResponse, error) {
	resp, err := restyClient.R().SetQueryParams(queryParams).Get("https://api.airvisual.com/v2/city")
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
	fmt.Println(string(resp.Body()))
	return body, nil
}
