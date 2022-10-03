package internal

import "net/url"

// creates authorization code request url from arguments
// response type and scope are openid in this project, because I did not implement other types and scopes LOL
// opEndpoint is authorization endpoint
func GetAuthCodeURL(opEndpoint, clientID, redirectURI, state, scope, responseType string) (string, error) {
	opURL, err := url.Parse(opEndpoint)
	if err != nil {
		return "", err
	}
	values := opURL.Query()
	values.Add("client_id", clientID)
	values.Add("redirect_uri", redirectURI)
	values.Add("state", state)
	values.Add("scope", scope)
	values.Add("response_type", responseType)
	opURL.RawQuery = values.Encode()
	return opURL.String(), nil
}
