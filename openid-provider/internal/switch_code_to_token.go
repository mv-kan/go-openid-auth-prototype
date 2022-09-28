package internal

import (
	"github.com/google/uuid"
	"github.com/mv-kan/go-openid-auth-prototype/internal/utils"
)

// Requires authRequest obj, deletes it and add generated token into token storage
func SwitchCodeToToken(authCode AuthCode) (Token, error) {
	// check authReq
	if !utils.ContainsID(AuthCodeStorage, authCode.GetID()) {
		return Token{}, utils.ErrNotFound
	}
	// add random id, scopes, and cliend id to token
	randID := uuid.New()
	accessToken := uuid.New()
	scopes := authCode.Scope
	clientID := authCode.ClientID
	// add token to token storage
	token := Token{ID: randID.String(), Scopes: scopes, ClientID: clientID, AccessToken: accessToken.String()}

	// remove auth req from storage and add token
	//  create new err value, cannot use := because it redifines RequestStorage
	var err error
	AuthCodeStorage, err = utils.RemoveByID(AuthCodeStorage, authCode.GetID())
	if err != nil {
		return token, err
	}
	TokenStorage = append(TokenStorage, token)

	// and return it
	return token, nil
}
