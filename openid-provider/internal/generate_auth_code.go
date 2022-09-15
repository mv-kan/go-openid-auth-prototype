package internal

import "github.com/google/uuid"

func GenerateAuthCode() AuthCode {
	randID := uuid.New().String()
	authCode := AuthCode(randID)
	AuthCodeStorage = append(AuthCodeStorage, authCode)
	return authCode
}
