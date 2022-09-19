package internal

import (
	"net/http"
	"net/url"
)

// errorCode - is OAuth Error code (https://www.rfc-editor.org/rfc/rfc6749#section-4.1.2), you can find then in pkg
func AuthErrorResponse(w http.ResponseWriter, r *http.Request, authReq AuthenticateRequest, errorCode string) {
	callback, err := url.Parse(authReq.RedirectURI)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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
