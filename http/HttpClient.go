package http

import (
	"bytes"
	"net/http"
	"time"
)

type HttpClient struct {
	client *http.Client
}

func NewHttpClient(timeout time.Duration) *HttpClient {
	return &HttpClient{
		client: &http.Client{Timeout: timeout},
	}
}

func (c *HttpClient) Get(url string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	addHeaders(req, headers)
	return c.client.Do(req)
}

func (c *HttpClient) Post(url string, body []byte, headers map[string]string) (*http.Response, error) {
	return c.doWithBody(http.MethodPost, url, body, headers)
}

func (c *HttpClient) Put(url string, body []byte, headers map[string]string) (*http.Response, error) {
	return c.doWithBody(http.MethodPut, url, body, headers)
}

func (c *HttpClient) Delete(url string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	addHeaders(req, headers)
	return c.client.Do(req)
}

func (c *HttpClient) doWithBody(method, url string, body []byte, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	addHeaders(req, headers)
	return c.client.Do(req)
}

func addHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}
