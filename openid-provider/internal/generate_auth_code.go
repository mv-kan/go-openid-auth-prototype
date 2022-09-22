package internal

import (
	"github.com/google/uuid"
	"github.com/mv-kan/go-openid-auth-prototype/internal/utils"
)

// add auth code obj to auth code storage and returns new created auth code
func GenerateAuthCode(authRequestID string) (AuthCode, error) {
	if !utils.ContainsID(RequestStorage, authRequestID) {
		return AuthCode{}, ErrAuthReqDoesNotExist
	}
	randID := uuid.New().String()

	authCode := AuthCode{ID: randID, AuthRequestID: authRequestID}
	AuthCodeStorage = append(AuthCodeStorage, authCode)
	return authCode, nil
}
