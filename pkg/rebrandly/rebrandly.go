package rebrandly

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cjtim/cjtim-backend-go/config"
	"github.com/cjtim/cjtim-backend-go/pkg/utils"
	"github.com/cjtim/cjtim-backend-go/repository"
)

type RebrandlyNewUrlReq struct {
	Destination string             `json:"destination"`
	Domain      RebrandlyDomainReq `json:"domain"`
}
type RebrandlyDomainReq struct {
	Fullname string `json:"fullName"`
}

// Add -
func Add(originalURL string) (*repository.URLScheama, error) {
	if originalURL[:8] != "https://" && originalURL[:7] != "http://" {
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
		"apikey":       config.Config.RebrandlyKey,
		"workspace":    config.Config.RebrandlyWorkspace,
	}
	resp, respBody, err := utils.Http(&utils.HttpReq{
		Method:  http.MethodPost,
		URL:     "https://api.rebrandly.com/v1/links",
		Headers: headers,
		Body:    body,
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(string(respBody))
	}
	data := repository.URLScheama{}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func Delete(id string) error {
	headers := map[string]string{
		"apikey": config.Config.RebrandlyKey,
	}
	resp, body, err := utils.Http(&utils.HttpReq{
		Method:  http.MethodDelete,
		URL:     "https://api.rebrandly.com/v1/links/" + id,
		Headers: headers,
	})
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(string(body))
	}
	return nil
}
