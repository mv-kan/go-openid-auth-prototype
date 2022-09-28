package internal

import (
	"github.com/google/uuid"
	"github.com/mv-kan/go-openid-auth-prototype/internal/utils"
)

// add auth code obj to auth code storage and returns new created auth code
func GenerateAuthCode(authReq AuthenticateRequest) (AuthCode, error) {
	if !utils.ContainsID(RequestStorage, authReq.GetID()) {
		return AuthCode{}, ErrAuthReqDoesNotExist
	}
	randID := uuid.New().String()

	authCode := AuthCode{ID: randID, AuthRequestID: authReq.GetID(), ClientID: authReq.ClientID, Scope: authReq.Scope, RedirectURI: authReq.RedirectURI, State: authReq.State}
	var err error
	RequestStorage, err = utils.RemoveByID(RequestStorage, authReq.GetID())
	if err != nil {
		return authCode, err
	}
	return authCode, nil
}
