package internal

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	IdToken     string `json:"id_token"`
	TokenType   string `json:"token_type"`
}
