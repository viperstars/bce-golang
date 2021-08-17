## 百度 BCE Golang SDK

##### 安装

```bash
go get github.com/viperstars/bce-golang
```

##### 使用

```go
package main

import (
	"fmt"
	"github.com/json-iterator/go"
    "github.com/viperstars/bce-golang/bce_client"
    "io/ioutil"
)

func main() {
	c := bce_client.NewBCEClient(
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
	fmt.Println(jsoniter.Get(b, "instances", 10).Get("id").ToString()) // use json-iterator or json.Unmashal
}

```