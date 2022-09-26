package internal

import (
	"github.com/google/uuid"
	"github.com/mv-kan/go-openid-auth-prototype/internal/utils"
)

// Requires authRequest obj, deletes it and add generated token into token storage
func GenerateAccessToken(authReq AuthenticateRequest) (Token, error) {
	// check authReq
	if !utils.ContainsID(RequestStorage, authReq.GetID()) {
		return Token{}, utils.ErrNotFound
	}
	// add random id, scopes, and cliend id to token
	randID := uuid.New()
	accessToken := uuid.New()
	scopes := authReq.Scope
	clientID := authReq.ClientID
	// add token to token storage
	token := Token{ID: randID.String(), Scopes: scopes, ClientID: clientID, AccessToken: accessToken.String()}

	// remove auth req from storage and add token
	utils.RemoveByID(RequestStorage, authReq.GetID())
	TokenStorage = append(TokenStorage, token)

	// and return it
	return token, nil
}
