package handler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/mv-kan/go-openid-auth-prototype/internal/utils"
	"github.com/mv-kan/go-openid-auth-prototype/internal/vars"
	"github.com/mv-kan/go-openid-auth-prototype/openid-provider/internal"
	"github.com/mv-kan/go-openid-auth-prototype/openid-provider/pkg"
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
	randID := uuid.New().String()
	authReqParams := internal.AuthenticateRequest{
		AuthRequestID: randID,
	}

	redirectURI, ok := params["redirect_uri"]
	if !ok {
		internal.AuthErrorResponse(w, r, authReqParams, pkg.InvalidRequest)
		return
	}
	authReqParams.RedirectURI = redirectURI[0]

	state, ok := params["state"]
	if !ok {
		internal.AuthErrorResponse(w, r, authReqParams, pkg.InvalidRequest)
		return
	}
	authReqParams.State = state[0]

	scope, ok := params["scope"]
	if !ok {
		internal.AuthErrorResponse(w, r, authReqParams, pkg.InvalidRequest)
		return
	}
	authReqParams.Scope = scope

	responseType, ok := params["response_type"]
	if !ok {
		internal.AuthErrorResponse(w, r, authReqParams, pkg.InvalidRequest)
		return
	}
	authReqParams.ResponseType = responseType

	clientID, ok := params["client_id"]
	if !ok {
		internal.AuthErrorResponse(w, r, authReqParams, pkg.InvalidRequest)
		return
	}
	authReqParams.ClientID = clientID[0]
	// verify that openid is present in scope
	if !utils.Contains(authReqParams.Scope, "openid") {
		internal.AuthErrorResponse(w, r, authReqParams, pkg.InvalidScope)
		return
	}
	// the only available flow is Authorization code flow
	if len(authReqParams.ResponseType) != 1 || authReqParams.ResponseType[0] != "code" {
		internal.AuthErrorResponse(w, r, authReqParams, pkg.UnsupportedResponseType)
		return
	}
	ids := func() []string {
		tmp := make([]string, 0, len(internal.ClientStorage))
		for _, client := range internal.ClientStorage {
			tmp = append(tmp, client.ID)
		}
		return tmp
	}()
	// error, client is not registered
	if !utils.Contains(ids, authReqParams.ClientID) {
		internal.AuthErrorResponse(w, r, authReqParams, pkg.UnauthorizedClient)
		return
	}
	internal.RequestStorage = append(internal.RequestStorage, authReqParams)

	// authenticate user redirecting him to login page
	loginRedirect, err := url.JoinPath(vars.OP_HOST, vars.LOGIN_ENDPOINT)
	loginRedirectParams := fmt.Sprintf("?authRequestID=%s", authReqParams.GetID())
	if err != nil {
		internal.AuthErrorResponse(w, r, authReqParams, pkg.ServerError)
		return
	}
	http.Redirect(w, r, loginRedirect+loginRedirectParams, http.StatusFound)
}

func authenticatePost(w http.ResponseWriter, r *http.Request) {
	return
}
