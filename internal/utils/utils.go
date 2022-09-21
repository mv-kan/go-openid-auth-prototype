package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrNotFound = fmt.Errorf("not found")
)

func ResponseJSON(w http.ResponseWriter, code int, payload any) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func GetByID[T IDer](sl []T, id string) (*T, error) {
	for i, value := range sl {
		if value.GetID() == id {
			return &sl[i], nil
		}
	}
	return nil, ErrNotFound
}

func ContainsID[T IDer](sl []T, id string) bool {
	for _, value := range sl {
		if value.GetID() == id {
			return true
		}
	}
	return false
}

func Contains[T comparable](sl []T, elem T) bool {
	for _, value := range sl {
		if value == elem {
			return true
		}
	}
	return false
}

func WriteResponse(w http.ResponseWriter, code int, message string) {
	w.Write([]byte(message))
	w.WriteHeader(code)
}

// Writes status code method is not allowed
// Also it writes in header all allowed methods
func AllowedMethods(w http.ResponseWriter, allowedMethods []string) {
	w.Header().Add("Allow", strings.Join(allowedMethods, ", "))
	w.WriteHeader(http.StatusMethodNotAllowed)
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
