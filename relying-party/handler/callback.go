package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/mv-kan/go-openid-auth-prototype/internal/log"
	"github.com/mv-kan/go-openid-auth-prototype/internal/utils"
	"github.com/mv-kan/go-openid-auth-prototype/internal/vars"
	"github.com/mv-kan/go-openid-auth-prototype/relying-party/config"
	"github.com/mv-kan/go-openid-auth-prototype/relying-party/internal"
)

func Callback(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	_, ok := params["code"]
	if ok {
		callbackToken(w, r)
	}
}

func callbackToken(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	codes, ok := params["code"]
	if !ok {
		log.Debug("parsing of code params failed")
		http.Error(w, "parsing of code params failed", http.StatusBadRequest)
		return
	}
	code := codes[0]
	states, ok := params["state"]
	if !ok {
		log.Debug("parsing of state params failed")
		http.Error(w, "parsing of state params failed", http.StatusBadRequest)
		return
	}
	state := states[0]
	cookieState, err := utils.GetCookie(r, "state")
	if err != nil {
		log.Debug(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if cookieState != state {
		log.Debug("states are not the same")
		http.Error(w, "states are not the same", http.StatusBadRequest)
		return
	}
	opEndpoint, err := url.JoinPath(vars.OP_URL, vars.TOKEN_ENDPOINT)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	// exchange code to token
	tokenRequest, err := internal.GetExchangeCodeToTokenRequest(opEndpoint, config.CLIENT_ID, config.CLIENT_SECRET, config.REDIRECT_URI, code, "authorization_code")
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	res, err := http.DefaultClient.Do(tokenRequest)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	// read all json body
	tokenResponse := internal.TokenResponse{}
	jsonBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	json.Unmarshal(jsonBytes, &tokenResponse)
	utils.SetCookie(w, "token", tokenResponse.AccessToken)
	http.Redirect(w, r, vars.RP_URL, http.StatusFound)
}
