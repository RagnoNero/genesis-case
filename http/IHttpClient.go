package http

import "net/http"

type IHttpClient interface {
	Get(url string, headers map[string]string) (*http.Response, error)
	Post(url string, body []byte, headers map[string]string) (*http.Response, error)
	Put(url string, body []byte, headers map[string]string) (*http.Response, error)
	Delete(url string, headers map[string]string) (*http.Response, error)
}
