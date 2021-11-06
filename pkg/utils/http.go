package utils

import (
	"io/ioutil"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func HttpGET(url string, querys, headers map[string]string) (resp *http.Response, body []byte, err error) {
	if querys == nil {
		querys = map[string]string{}
	}
	if headers == nil {
		headers = map[string]string{}
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return &http.Response{}, nil, err
	}

	// appending querys
	q := req.URL.Query()
	for query, value := range querys {
		q.Add(query, value)
	}
	req.URL.RawQuery = q.Encode()

	// appending headers
	for header, value := range headers {
		req.Header.Add(header, value)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err = client.Do(req)
	if err != nil {
		return &http.Response{}, nil, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("HttpGET", zap.Error(err))
	}

	return &http.Response{}, body, nil
}
