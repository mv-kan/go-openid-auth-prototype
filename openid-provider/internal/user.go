package internal

import (
	"fmt"

	"github.com/mv-kan/go-openid-auth-prototype/internal/utils"
)

type User struct {
	Username string
	Password string
}

func CheckUsernamePassword(username, password, authReqID string) error {
	if !utils.ContainsID(RequestStorage, authReqID) {
		return fmt.Errorf("there is no such authorization request")
	}

	//for demonstration purposes we'll check on a static list with plain text password
	//for real world scenarios, be sure to have the password hashed and salted (e.g. using bcutils
	for _, user := range UserStorage {
		if user.Username == username && user.Password == password {
			return nil
		}
	}
	return fmt.Errorf("username or password wrong")
}
