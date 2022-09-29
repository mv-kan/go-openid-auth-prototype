package handler

import (
	"fmt"
	"net/http"

	"github.com/mv-kan/go-openid-auth-prototype/internal/log"
	"github.com/mv-kan/go-openid-auth-prototype/internal/utils"
	"github.com/mv-kan/go-openid-auth-prototype/openid-provider/internal"
)

func CheckToken(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		checkTokenPost(w, r)
		return
	default:
		utils.AllowedMethods(w, []string{http.MethodPost})
		return
	}
}

// requires token key value in a post body
func checkTokenPost(w http.ResponseWriter, r *http.Request) {
	// authorization code check
	err := r.ParseForm()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, fmt.Sprintf("cannot parse form:%s", err), http.StatusInternalServerError)
		return
	}
	accessToken := r.FormValue("token")
	// check token by getting it from "db"
	token := getTokenByAccessToken(internal.TokenStorage, accessToken)
	if token == nil {
		log.Debug("token is not present in token storage")
		http.Error(w, "no such token", http.StatusUnauthorized)
		return
	}
	// yes I do understand that we don't use scopes at all
	// but this project is for example and learning reasons
	w.WriteHeader(http.StatusOK)
}
func getTokenByAccessToken(tokenStorage []internal.Token, accessToken string) *internal.Token {
	for _, val := range tokenStorage {
		if val.AccessToken == accessToken {
			return &val
		}
	}
	return nil
}
