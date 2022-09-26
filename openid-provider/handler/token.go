package handler

import (
	"crypto/sha256"
	"crypto/subtle"
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
		return
	}
	// Check Authorization basic Client
	client, err := utils.GetByID(internal.ClientStorage, clientID)
	if errors.Is(err, utils.ErrNotFound) {
		log.Debug(fmt.Sprintf("clientID does not exist clientID=%s", clientID))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Error(err.Error())
		http.Error(w, fmt.Sprintf("cannot get auth request:%s", err), http.StatusInternalServerError)
		return
	}

	clientID, clientSecret, ok := r.BasicAuth()
	if !ok {
		log.Debug("no basic auth is presented")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	clientIDHash := sha256.Sum256([]byte(clientID))
	secretHash := sha256.Sum256([]byte(clientSecret))
	expectedClientIDHash := sha256.Sum256([]byte(client.ID))
	expectedSecretHash := sha256.Sum256([]byte(client.Secret))

	// Use the subtle.ConstantTimeCompare() function to check if
	// the provided username and password hashes equal the
	// expected username and password hashes. ConstantTimeCompare
	// will return 1 if the values are equal, or 0 otherwise.
	// Importantly, we should to do the work to evaluate both the
	// username and password before checking the return values to
	// avoid leaking information.
	clientIDMatch := (subtle.ConstantTimeCompare(clientIDHash[:], expectedClientIDHash[:]) == 1)
	secretMatch := (subtle.ConstantTimeCompare(secretHash[:], expectedSecretHash[:]) == 1)
	if !(clientIDMatch && secretMatch) {
		log.Debug(fmt.Sprintf("clientid or password is invalid, clientID=%s", clientID))
		http.Error(w, "clientid or password is invalid, unauthorized", http.StatusUnauthorized)
		return
	}

	// after successful validation we generate token and add it to token storage
	token, err := internal.GenerateAccessToken(*authReq)
	if errors.Is(err, utils.ErrNotFound) {
		log.Debug(fmt.Sprintf("clientID does not exist clientID=%s", clientID))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Error(err.Error())
		http.Error(w, fmt.Sprintf("cannot generate access token:%s", err), http.StatusInternalServerError)
		return
	}
	// then write successful response
	w.Header().Add("Pragma", "no-cache")
	w.Header().Add("Cache-Control", "no-store")

	msg := map[string]any{
		"access_token": token.AccessToken,
		"token_type":   "Bearer",
		"expires_in":   99999999,
		"id_token":     token.GetID(),
	}
	err = utils.ResponseJSON(w, http.StatusOK, msg)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, fmt.Sprintf("cannot generate access token:%s", err), http.StatusInternalServerError)
	}
}
