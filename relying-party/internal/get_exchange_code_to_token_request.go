package internal

import (
	"net/http"
	"net/url"
	"strings"
)

// grant_type = authorization_code this is always because I didnt implmenet other grant types in openid provider
// opEndpoint is token exchange endpoint
func GetExchangeCodeToTokenRequest(opEndpoint, clientID, clientSecret, redirectURI, code, grantType string) (*http.Request, error) {
	data := url.Values{
		"code":         {code},
		"redirect_uri": {redirectURI},
		"client_id":    {clientID},
		"grant_type":   {grantType},
	}
	request, err := http.NewRequest(http.MethodPost, opEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(clientID, clientSecret)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return request, nil
}
