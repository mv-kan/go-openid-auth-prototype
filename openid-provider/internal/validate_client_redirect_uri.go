package internal

import (
	"errors"
	"net/url"

	"github.com/mv-kan/go-openid-auth-prototype/internal/utils"
)

// returns true if redirectURI is valid for clientID, otherwise false
func ValidateClientRedirectURI(clientID string, redirectURI string) (bool, error) {
	// get client
	client, err := utils.GetByID(ClientStorage, clientID)
	if err != nil {
		return false, errors.New("client does not exist")
	}
	// check with ParseRequestURI all urls for validity and then compare them
	_, err = url.ParseRequestURI(redirectURI)
	if err != nil {
		return false, err
	}
	// check if redirect uri is in client registered redirect uris list
	return utils.Contains(client.RedirectURIs, redirectURI), nil
}
