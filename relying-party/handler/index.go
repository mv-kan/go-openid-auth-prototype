package handler

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/mv-kan/go-openid-auth-prototype/internal/log"
	"github.com/mv-kan/go-openid-auth-prototype/internal/utils"
	"github.com/mv-kan/go-openid-auth-prototype/internal/vars"
	"github.com/mv-kan/go-openid-auth-prototype/relying-party/config"
	"github.com/mv-kan/go-openid-auth-prototype/relying-party/internal"
)

func Index(w http.ResponseWriter, r *http.Request) {
	token, err := utils.GetCookie(r, "token")
	// if token does not exist, then redirect to auth server's login page
	if errors.Is(err, utils.ErrNotFound) {
		// save state in cookies
		state := uuid.New().String()
		utils.SetCookie(w, "state", state)
		// get openid provider authenticate url
		opEndpoint, err := url.JoinPath(vars.OP_URL, vars.AUTHN_ENDPOINT)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		redirectURL, err := internal.GetAuthCodeURL(opEndpoint, config.CLIENT_ID, config.REDIRECT_URI, state, "openid", "code")
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, redirectURL, http.StatusFound)
	} else if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeProtectedInfo(token, w, r)
}
func writeProtectedInfo(token string, w http.ResponseWriter, r *http.Request) {

	// if token exists send token to protected resource and send message from this resource
	url, err := url.JoinPath(vars.RS_URL, vars.RS_PROTECTED_ENDPOINT)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + token
	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responseBody := []byte("Message from protected server: ")
	responseBody = append(responseBody, body...)
	w.Write(responseBody)
}
