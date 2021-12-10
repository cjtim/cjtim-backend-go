package rebrandly

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cjtim/cjtim-backend-go/configs"
	"github.com/cjtim/cjtim-backend-go/internal/app/repository"
	"github.com/cjtim/cjtim-backend-go/internal/pkg/utils"
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
	noProtocol := originalURL[:8] != "https://" && originalURL[:7] != "http://"
	if noProtocol {
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
		"apikey":       configs.Config.RebrandlyKey,
		"workspace":    configs.Config.RebrandlyWorkspace,
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
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(respBody))
	}
	data := repository.URLScheama{}
	err = json.Unmarshal(respBody, &data)
	return &data, err
}

func Delete(id string) error {
	headers := map[string]string{
		"apikey": configs.Config.RebrandlyKey,
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
