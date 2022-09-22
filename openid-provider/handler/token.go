package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mv-kan/go-openid-auth-prototype/internal/log"
	"github.com/mv-kan/go-openid-auth-prototype/internal/utils"
	"github.com/mv-kan/go-openid-auth-prototype/openid-provider/internal"
)

func Token(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		tokenPost(w, r)
		return
	default:
		utils.AllowedMethods(w, []string{http.MethodPost})
		return
	}
}

func tokenPost(w http.ResponseWriter, r *http.Request) {
	// authorization code check
	err := r.ParseForm()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, fmt.Sprintf("cannot parse form:%s", err), http.StatusInternalServerError)
		return
	}
	code := r.FormValue("code")
	redirectURI := r.FormValue("redirect_uri")
	clientID := r.FormValue("client_id")
	grantType := r.FormValue("grant_type")
	// validate grantType
	if grantType != "authorization_code" {
		log.Debug("only grant type is authorization_code, grant_type = " + grantType)
		http.Error(w, "the only grant type is authorization_code, grant_type="+grantType, http.StatusBadRequest)
		return
	}
	// get auth code obj from storage
	authCode, err := utils.GetByID(internal.AuthCodeStorage, code)
	if errors.Is(err, utils.ErrNotFound) {
		log.Debug(err.Error())
		http.Error(w, fmt.Sprintf("cannot get auth code:%s", err), http.StatusBadRequest)
		return
	} else if err != nil {
		log.Error(err.Error())
		http.Error(w, fmt.Sprintf("cannot get auth code:%s", err), http.StatusInternalServerError)
		return
	}
	// get auth req obj from storage using auth code
	authReq, err := utils.GetByID(internal.RequestStorage, authCode.AuthRequestID)
	if errors.Is(err, utils.ErrNotFound) {
		log.Debug(err.Error())
		http.Error(w, fmt.Sprintf("cannot get auth request:%s", err), http.StatusBadRequest)
		return
	} else if err != nil {
		log.Error(err.Error())
		http.Error(w, fmt.Sprintf("cannot get auth request:%s", err), http.StatusInternalServerError)
		return
	}
	// compare redirect uri
	if authReq.RedirectURI != redirectURI {
		log.Debug(fmt.Sprintf("redirect uris are not the same: %s(in auth req) %s (in token req)", authReq.RedirectURI, redirectURI))
		http.Error(w, "redirect uris are not the same", http.StatusBadRequest)
		return
	}
	// validate clientID
	if authReq.ClientID != clientID {
		log.Debug(fmt.Sprintf("clientIDs are not the same clientid auth req = %s, clientid token req = %s", authReq.ClientID, clientID))
		http.Error(w, "client ids are not the same", http.StatusBadRequest)
	}
}
