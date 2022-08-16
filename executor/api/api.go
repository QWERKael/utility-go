package api

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// 表示请求未成功的状态码
var errStatusCode = 0

func Get(url string, headers map[string]string, urlQuery map[string]string) (int, []byte, error) {
	return Request(http.MethodGet, url, headers, urlQuery, nil)
}

func Post(url string, headers map[string]string, urlQuery map[string]string, data []byte) (int, []byte, error) {
	return Request(http.MethodPost, url, headers, urlQuery, data)
}

func Put(url string, headers map[string]string, urlQuery map[string]string, data []byte) (int, []byte, error) {
	return Request(http.MethodPut, url, headers, urlQuery, data)
}

func Request(method string, url string, headers map[string]string, urlQuery map[string]string, data []byte) (int, []byte, error) {
	client := &http.Client{}
	// 创建新的 request
	req, err := http.NewRequest(method, url, strings.NewReader(string(data)))
	if err != nil {
		return errStatusCode, nil, err
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
	// 执行 request
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return errStatusCode, nil, err
	}
	defer resp.Body.Close()
	// 返回结果
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return errStatusCode, nil, err
	}
	return resp.StatusCode, body, nil
}
