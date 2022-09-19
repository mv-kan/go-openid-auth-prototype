package internal

import (
	"net/http"
	"net/url"

	"github.com/mv-kan/go-openid-auth-prototype/internal/utils"
)

// returns json object with errorCode, also it has redirection uri with error params
// errorCode - is OAuth Error code (https://www.rfc-editor.org/rfc/rfc6749#section-4.1.2), you can find then in pkg
func AuthErrorResponse(w http.ResponseWriter, r *http.Request, authReq AuthenticateRequest, errorCode string) {
	callback, err := url.Parse(authReq.RedirectURI)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// auth server MUST NOT redirect if redirectURI is invalid
	ok, err := ValidateClientRedirectURI(authReq.ClientID, callback.String())
	if err.Error() == "client does not exist" {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "client does not exist"})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "redirectURI is not valid. Maybe it does not exist in client registered redirectURIs list"})
		return
	}
	// add to callback query parameters
	values := callback.Query()

	// add state parameter if exists
	if authReq.State != "" {
		values.Add("state", authReq.State)
	}

	// add error parameter according to https://www.rfc-editor.org/rfc/rfc6749#section-4.1.2
	values.Add("error", errorCode)

	callback.RawQuery = values.Encode()

	// oauth specification requires add form urlencoded
	w.Header().Add("Content-Type", "application/x-www-form-urlencoded")
	http.Redirect(w, r, callback.String(), http.StatusFound)
}
