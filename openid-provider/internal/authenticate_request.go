package internal

import "net/url"

type AuthenticateRequest struct {
	Scope         []string
	ResponseType  []string
	ClientID      string
	RedirectURI   string
	State         string
	AuthRequestID string
}

func (a AuthenticateRequest) GetID() string {
	return a.AuthRequestID
}

func (a AuthenticateRequest) GetCallbackURL(authCode AuthCode) (string, error) {
	redirectURL, err := url.Parse(a.RedirectURI)
	if err != nil {
		return "", err
	}

	values := redirectURL.Query()
	values.Add("code", authCode.ID)
	values.Add("state", a.State)
	redirectURL.RawQuery = values.Encode()
	return redirectURL.String(), nil
}

func (a AuthenticateRequest) GetCallbackURLAuto() (string, error) {
	authCode, err := GenerateAuthCode(a.GetID())
	if err != nil {
		return "", err
	}
	url, err := a.GetCallbackURL(authCode)
	return url, err
}
