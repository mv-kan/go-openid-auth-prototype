package handler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/mv-kan/go-openid-auth-prototype/internal/utils"
	"github.com/mv-kan/go-openid-auth-prototype/internal/vars"
	"github.com/mv-kan/go-openid-auth-prototype/openid-provider/internal"
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
	params := r.URL.Query()
	// parse parameters
	scope, ok := params["scope"]
	if !ok {
		return
	}
	responseType, ok := params["response_type"]
	if !ok {
		return
	}
	clientID, ok := params["client_id"]
	if !ok {
		return
	}
	redirectURI, ok := params["redirect_uri"]
	if !ok {
		return
	}
	state, ok := params["state"]
	if !ok {
		return
	}

	randID := uuid.New().String()
	authReqParams := internal.AuthenticateRequest{
		AuthRequestID: randID,
		Scope:         scope,
		ResponseType:  responseType,
		ClientID:      clientID[0],
		RedirectURI:   redirectURI[0],
		State:         state[0],
	}

	// verify that openid is present in scope
	if !utils.Contains(authReqParams.Scope, "openid") {
		return
	}
	if len(authReqParams.ResponseType) != 1 || authReqParams.ResponseType[0] != "code" {
		return
	}
	ids := func() []string {
		tmp := make([]string, 0, len(internal.ClientStorage))
		for _, client := range internal.ClientStorage {
			tmp = append(tmp, client.ID)
		}
		return tmp
	}()
	if !utils.Contains(ids, authReqParams.ClientID) {
		return
	}
	internal.RequestStorage = append(internal.RequestStorage, authReqParams)

	// authenticate user redirecting him to login page
	loginRedirect, err := url.JoinPath(vars.AUTH_SERVER, vars.LOGIN_ENDPOINT, fmt.Sprintf("?authRequestID=%s", randID))
	if err != nil {
		return
	}
	http.Redirect(w, r, loginRedirect, http.StatusFound)
}

func authenticatePost(w http.ResponseWriter, r *http.Request) {
	return
}
