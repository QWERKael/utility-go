package api

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

// 表示请求未成功的状态码
var errStatusCode = 0

func Get(url string, headers map[string]string, urlQuery map[string]string) (int, []byte, error) {
	code, body, _, err := Request(http.MethodGet, url, headers, urlQuery, nil)
	return code, body, err
}

func Post(url string, headers map[string]string, urlQuery map[string]string, data []byte) (int, []byte, error) {
	code, body, _, err := Request(http.MethodPost, url, headers, urlQuery, data)
	return code, body, err
}

func Put(url string, headers map[string]string, urlQuery map[string]string, data []byte) (int, []byte, error) {
	code, body, _, err := Request(http.MethodPut, url, headers, urlQuery, data)
	return code, body, err
}

func ReTryRequest(f func(args ...interface{}) (int, []byte, error), reTryTimes int) func(...interface{}) (int, []byte, error) {
	return func(args ...interface{}) (int, []byte, error) {
		for i := 0; i < reTryTimes; i++ {
			statusCode, rst, err := f(args)
			if err != nil {
				return errStatusCode, nil, err
			}
			if statusCode != 200 {
				continue
			}
			return statusCode, rst, nil
		}
		return errStatusCode, nil, errors.New("重试失败")
	}
}

func Request(method string, url string, headers map[string]string, urlQuery map[string]string, data []byte) (int, []byte, string, error) {
	client := &http.Client{}
	// 创建新的 request
	req, err := http.NewRequest(method, url, strings.NewReader(string(data)))
	if err != nil {
		return errStatusCode, nil, "", err
	}
	// 添加 header
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	// 添加 url query
	q := req.URL.Query()
	for queryKey, queryValue := range urlQuery {
		q.Set(queryKey, queryValue)
	}
	req.URL.RawQuery = q.Encode()
	uri := req.RequestURI
	// 执行 request
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return errStatusCode, nil, "", err
	}
	defer resp.Body.Close()
	// 返回结果
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return errStatusCode, nil, "", err
	}
	return resp.StatusCode, body, uri, nil
}
