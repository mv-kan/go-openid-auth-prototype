package utils

import (
	"net/http"
	"strings"
)

func WriteResponse(w http.ResponseWriter, code int, message string) {
	w.Write([]byte(message))
	w.WriteHeader(code)
}

func MakeRequest(method string, url string, body string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return res, err
}
