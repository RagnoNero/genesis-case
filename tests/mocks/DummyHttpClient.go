package mocks

import (
	"net/http"
)

type DummyHttpClient struct {
	ReturnFunc func() (*http.Response, error)
}

func NewDummyHttp(retFunc func() (*http.Response, error)) *DummyHttpClient {
	return &DummyHttpClient{
		ReturnFunc: retFunc,
	}
}

func (c *DummyHttpClient) Get(url string, headers map[string]string) (*http.Response, error) {
	return c.ReturnFunc()
}

func (c *DummyHttpClient) Post(url string, body []byte, headers map[string]string) (*http.Response, error) {
	return &http.Response{}, nil
}

func (c *DummyHttpClient) Put(url string, body []byte, headers map[string]string) (*http.Response, error) {
	return &http.Response{}, nil
}

func (c *DummyHttpClient) Delete(url string, headers map[string]string) (*http.Response, error) {
	return &http.Response{}, nil
}
