package utils

import (
	"errors"
	"net/http"
)

func SetCookie(w http.ResponseWriter, name, value string) {
	cookie := &http.Cookie{
		Name:  name,
		Value: value,
	}
	http.SetCookie(w, cookie)
}

// ErrNotFound if cookie is not found
func GetCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if errors.Is(err, http.ErrNoCookie) {
		return "", ErrNotFound
	} else if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
