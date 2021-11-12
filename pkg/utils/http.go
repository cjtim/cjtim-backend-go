package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpReq struct {
	Method    string
	URL       string
	Querys    map[string]string
	Headers   map[string]string
	Body      interface{}
	BodyBytes []byte
}

func doRequest(req *http.Request) (resp *http.Response, body []byte, err error) {
	client := &http.Client{Timeout: 15 * time.Second}
	fmt.Printf("req.RequestURI: %v\n", req.URL.RawQuery)
	resp, err = client.Do(req)
	if err != nil {
		return &http.Response{}, nil, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return &http.Response{}, nil, err
	}
	return resp, body, nil
}

func Http(httpReq *HttpReq) (*http.Response, []byte, error) {
	if httpReq.Querys == nil {
		httpReq.Querys = map[string]string{}
	}
	if httpReq.Headers == nil {
		httpReq.Headers = map[string]string{}
	}
	if httpReq.Body != nil {
		bBody, err := json.Marshal(httpReq.Body)
		httpReq.BodyBytes = bBody
		if err != nil {
			return nil, nil, err
		}
	}

	req, err := http.NewRequest(
		httpReq.Method,
		httpReq.URL,
		bytes.NewBuffer(httpReq.BodyBytes),
	)
	if err != nil {
		return nil, nil, err
	}

	// appending querys
	rawQuery := ""
	for query, value := range httpReq.Querys {
		rawQuery += fmt.Sprintf("&%s=%s", query, url.QueryEscape(value))
	}
	req.URL.RawQuery = strings.TrimPrefix(rawQuery, "&")

	// appending headers
	for header, value := range httpReq.Headers {
		req.Header.Add(header, value)
	}
	return doRequest(req)
}
