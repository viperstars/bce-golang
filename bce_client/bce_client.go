package bce_client

import (
	"bytes"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/satori/go.uuid"
	"github.com/viperstars/bce-golang/bce_http"
	"github.com/viperstars/bce-golang/bce_signer"
	"io/ioutil"
	"net/http"
)

type client struct {
	accessKey string
	secretKey string
	region    string
	instance  string
	host      string
}

type BCEClient struct {
	*client
}

func (c *client) Execute(method string, url string, params map[string]string, headers map[string]string, body []byte) (*http.Response, string, error) {
	if _, ok := params["clientToken"]; !ok {
		clientToken := uuid.NewV4()
		params["clientToken"] = clientToken.String()
	}
	req := &bce_http.Request{
		BaseUrl: fmt.Sprintf("http://%s", c.host),
		Method:  method,
		Path:    url,
		Query:   params,
		Headers: headers,
		Body:    bytes.NewReader(body),
	}
	return c.doRequest(req)
}

func NewBCEClient(accessKey string, secretKey string, region string, instance string) *BCEClient {
	host := fmt.Sprintf("%s.%s.baidubce.com", instance, region)
	client := &client{accessKey, secretKey, region, instance, host}
	return &BCEClient{client}
}

func (c *client) doRequest(req *bce_http.Request) (*http.Response, string, error) {

	req.Headers[bce_http.HOST] = c.host

	timestamp := bce_signer.GetHttpHeadTimeStamp()
	authorization := bce_signer.Sign(c.accessKey, c.secretKey, timestamp, req.Method, req.Path, req.Query, req.Headers)

	req.Headers[bce_http.DATE] = timestamp
	req.Headers[bce_http.AUTHORIZATION] = authorization

	res, err := bce_http.Run(req, nil)
	return res, req.Query["clientToken"], err
}

func example() {
	c := NewBCEClient(
		"access",
		"secret",
		"bj",
		"bcc",
	)
	param := make(map[string]string)
	param["marker"] = "i-h7VoqzIT"
	headers := make(map[string]string)
	resp, token, err := c.Execute("GET", "/v2/instance", param, headers, nil)
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(token, err)
	fmt.Println(jsoniter.Get(b, "instances", 10).Get("id").ToString())
}
