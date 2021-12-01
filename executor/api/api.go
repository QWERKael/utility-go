package api

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, nil
}

func Put(url string, data []byte) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
