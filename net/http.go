package net

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func PostWithHeader(url string, msg []byte, headers map[string]string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(msg)))
	if err != nil {
		return nil, err
	}
	for key, header := range headers {
		req.Header.Set(key, header)
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
