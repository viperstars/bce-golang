package bce_signer

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

const BCEPREFIX = "x-bce-"

var Debug bool

func getCanonicalHeaders(headers map[string]string, headersToSign []string) (string, string) {
	if headersToSign == nil {
		headersToSign = []string{"host", "content-md5", "content-length", "content-type"}
	}

	result := make([]string, 0)
	signedHeaders := make([]string, 0)
	for k, v := range headers {
		k = strings.ToLower(k)
		if strings.HasPrefix(k, BCEPREFIX) || IsStringInSlice(k, headersToSign) {
			key := UriEncode(k)
			value := UriEncode(v)
			strTemplate := fmt.Sprintf("%s:%s", key, value)
			result = append(result, strTemplate)
			signedHeaders = append(signedHeaders, fmt.Sprint(key))
		}
	}
	sort.Strings(result)

	return strings.Join(result, "\n"), strings.Join(signedHeaders, ";")
}

func getCannonicalQuery(query map[string]string) string {
	if len(query) == 0 {
		return ""
	}

	result := make([]string, 0)
	for k, v := range query {
		if v != "" {
			key := UriEncode(k)
			value := UriEncode(v)
			result = append(result, fmt.Sprintf("%s=%s", key, value))
		} else {
			key := UriEncode(k)
			result = append(result, fmt.Sprintf("%s=", key))
		}
	}
	sort.Strings(result)

	return strings.Join(result, "&")
}

func Sign(accessKey, secretKey, timestamp, httpMethod, path string, query map[string]string,
	headers map[string]string) string {

	if path[0] != '/' {
		path = "/" + path
	}

	var expirationPeriodInSeconds = 1800
	authStringPrefix := fmt.Sprintf("bce-auth-v1/%s/%s/%d", accessKey,
		timestamp, expirationPeriodInSeconds)
	//fmt.Println(authStringPrefix)

	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(authStringPrefix))
	signingKey := fmt.Sprintf("%x", mac.Sum(nil))
	//fmt.Printf(signingKey)

	CanonicalURI := UriEncodeExceptSlash(path)
	CanonicalQueryString := getCannonicalQuery(query)
	fmt.Println(CanonicalQueryString)
	CanonicalHeaders, signedHeaders := getCanonicalHeaders(headers, nil)
	CanonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s", httpMethod, CanonicalURI,
		CanonicalQueryString, CanonicalHeaders)

	mac = hmac.New(sha256.New, []byte(signingKey))
	mac.Write([]byte(CanonicalRequest))
	signature := fmt.Sprintf("%x", mac.Sum(nil))

	authorization := fmt.Sprintf("%s/%s/%s", authStringPrefix, signedHeaders, signature)
	if Debug {
		fmt.Println(CanonicalRequest)
		fmt.Println(authorization)
	}
	return authorization
}

func GetHttpHeadTimeStamp() string {
	gmt := time.Now().UTC()
	return gmt.Format("2006-01-02T15:04:05Z")
}

func IsStringInSlice(s string, slice []string) bool {
	for _, v := range slice {
		if s == v {
			return true
		}
	}
	return false
}

func UriEncode(s string) string {
	result := make([]string, 0)
	ss := strings.Split(s, "/")
	if len(ss) == 1 {
		value := url.QueryEscape(s)
		return value
	}
	for _, v := range ss {
		value := url.QueryEscape(v)
		result = append(result, value)
	}
	return strings.Join(result, "%2F")
}

func UriEncodeExceptSlash(s string) string {
	result := make([]string, 0)
	ss := strings.Split(s, "/")
	if len(ss) == 1 {
		value := url.QueryEscape(s)
		return value
	}
	for _, v := range ss {
		value := url.QueryEscape(v)
		result = append(result, value)
	}
	return strings.Join(result, "/")
}
