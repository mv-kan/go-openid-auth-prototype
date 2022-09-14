package handler

import (
	"net/http"

	"github.com/mv-kan/go-openid-auth-prototype/internal/utils"
)

func Authenticate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		authenticateGet(w, r)
		return
	case http.MethodPost:

		return
	default:
		utils.AllowedMethods(w, []string{http.MethodGet, http.MethodPost})
		return
	}
}

func authenticateGet(w http.ResponseWriter, r *http.Request) {
	r.URL.Query()
}

func authenticatePost(w http.ResponseWriter, r *http.Request) {
	return
}
