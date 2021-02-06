package rebrandly

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/go-resty/resty/v2"
)

type RebrandlyNewUrlReq struct {
	Destination string             `json:"destination"`
	Domain      RebrandlyDomainReq `json:"domain"`
}
type RebrandlyDomainReq struct {
	Fullname string `json:"fullName"`
}
type RebrandlyNewUrlResp struct {
	ID          string `json:"id"`
	ShortURL    string `json:"shortUrl"`
	Destination string `json:"destination"`
}

var restyClient = resty.New()

// Add -
func Add(originalURL string) (*RebrandlyNewUrlResp, error) {
	if originalURL[:8] != "https://" || originalURL[:7] != "http://" {
		originalURL = "http://" + originalURL
	}
	body := &RebrandlyNewUrlReq{
		Destination: originalURL,
		Domain: RebrandlyDomainReq{
			Fullname: "link.cjtim.com",
		},
	}
	headers := map[string]string{
		"Content-Type": "application/json",
		"apikey":       os.Getenv("REBRANDLY_API"),
		"workspace":    os.Getenv("REBRANDLY_WORDSPACE"),
	}
	resp, err := restyClient.R().SetHeaders(headers).SetBody(body).Post("https://api.rebrandly.com/v1/links")
	if err != nil {
		return nil, err
	}
	data := &RebrandlyNewUrlResp{}
	if err := json.Unmarshal(resp.Body(), data); err != nil {
		return nil, err
	}
	return data, nil
}

func Delete(id string) error {
	resp, err := restyClient.R().Delete("https://api.rebrandly.com/v1/links/" + id)
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return errors.New(string(resp.Body()))
	}
	return nil
}
