package bce_http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var Debug bool = false

const (
	AUTHORIZATION       = "Authorization"
	CONTENT_DISPOSITION = "Content-Disposition"
	CONTENT_ENCODING    = "Content-Encoding"
	CONTENT_LENGTH      = "Content-Length"
	CONTENT_MD5         = "Content-MD5"
	CONTENT_RANGE       = "Content-Range"
	CONTENT_TYPE        = "Content-Type"
	DATE                = "Date"
	ETAG                = "ETag"
	EXPIRES             = "Expires"
	HOST                = "Host"
	LAST_MODIFIED       = "Last-Modified"
	RANGE               = "Range"
	SERVER              = "Server"
	USER_AGENT          = "User-Agent"
	OCTET_STREAM        = "application/octet-stream"
)

type Request struct {
	Method  string
	Path    string
	Query   map[string]string
	Headers map[string]string
	BaseUrl string
	Type    string
	Body    *bytes.Reader
	Timeout time.Duration
}

func (req *Request) url() (*url.URL, error) {
	u, err := url.Parse(req.BaseUrl)
	if err != nil {
		return nil, fmt.Errorf("bad endpoint url %q: %v", req.BaseUrl, err)
	}
	u.RawQuery = getUrl(req.Query)
	u.Path = req.Path
	return u, nil
}

func initHttpRequest(req *Request) (*http.Request, error) {
	requestUrl, _ := req.url()
	newReq, _ := http.NewRequest(req.Method, requestUrl.String(), nil)
	for k, v := range req.Headers {
		newReq.Header.Add(k, v)
	}
	if req.Body != nil {
		newReq.Body = ioutil.NopCloser(req.Body)
		newReq.ContentLength = int64(req.Body.Len())
		newReq.Header.Add(CONTENT_LENGTH, fmt.Sprintf("%d", newReq.ContentLength))
		if req.Type != "" {
			newReq.Header.Add(CONTENT_TYPE, req.Type)
		} else {
			newReq.Header.Add(CONTENT_TYPE, OCTET_STREAM)
		}
	}
	return newReq, nil
}

func doHttpRequest(httpClient *http.Client, req *http.Request, res interface{}) (*http.Response, error) {
	response, err := httpClient.Do(req)
	if Debug {
		fmt.Println("request: ", req)
		fmt.Println("response: ", response)
		fmt.Println("error: ", err)
	}
	return response, err
}

func Run(req *Request, res interface{}) (*http.Response, error) {
	hreq, _ := initHttpRequest(req)
	httpClient := &http.Client{
		Timeout: req.Timeout,
	}
	return doHttpRequest(httpClient, hreq, res)
}

func getUrl(q map[string]string) string {
	encodedUrl := &url.Values{}
	for k, v := range q {
		encodedUrl.Add(k, v)
	}
	return encodedUrl.Encode()
}
