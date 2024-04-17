package http_utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func GetRequestClosure(reqUrl string, headers map[string]string, params map[string]string) (func() (*http.Response, error), error) {
	client := &http.Client{}

	if params != nil {
		urlParams := url.Values{}
		for k, v := range params {
			urlParams.Add(k, v)
		}
		reqUrl = fmt.Sprintf("%s?%s", reqUrl, urlParams.Encode())
	}

	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		for header, val := range headers {
			req.Header.Set(header, val)
		}
	}
	return func() (*http.Response, error) {
		return client.Do(req)
	}, nil
}

func GetRequest(reqUrl string, headers map[string]string, params map[string]string) (*http.Response, error) {
	client := &http.Client{}

	if params != nil {
		urlParams := url.Values{}
		for k, v := range params {
			urlParams.Add(k, v)
		}
		reqUrl = fmt.Sprintf("%s?%s", reqUrl, urlParams.Encode())
	}

	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		for header, val := range headers {
			req.Header.Set(header, val)
		}
	}
	return client.Do(req)
}

func PostRequest(url string, headers map[string]string, body map[string]any) (*http.Response, error) {
	client := &http.Client{}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	if headers != nil {
		for header, val := range headers {
			req.Header.Set(header, val)
		}
	}

	return client.Do(req)
}
