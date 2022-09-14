package middleware

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/mv-kan/go-openid-auth-prototype/resource-server/utils"
)

var (
	authServer         = "http://localhost:7001"
	checkTokenEndpoint = "check-token"
)

func OnlyMethod(method string, allowedMethods []string, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == method {
			f(w, r)
		} else {
			w.Header().Add("Allow", strings.Join(allowedMethods, ", "))
			utils.WriteResponse(w, http.StatusMethodNotAllowed, "Not allowed")
		}
	}
}

func GetProtectedSuperSecret(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi, it is sure really cool access token you have here"))
}

func ValidateToken(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")

		// get token without bearer
		tmp := strings.Split(bearerToken, "Bearer ")
		if len(tmp) != 2 {
			utils.WriteResponse(w, http.StatusUnauthorized, "invalid authorization header")
			return
		}
		token := tmp[1]

		// check token
		// send token to server
		checkTokenURL, err := url.JoinPath(authServer, checkTokenEndpoint)
		if err != nil {
			utils.WriteResponse(w, http.StatusInternalServerError, "internal error")
		}

		// get response
		res, err := utils.MakeRequest(http.MethodPost, checkTokenURL, token)
		if err != nil {
			utils.WriteResponse(w, http.StatusInternalServerError, "durring request error occured")
			return
		}
		if res.StatusCode != http.StatusOK {
			utils.WriteResponse(w, http.StatusUnauthorized, "access token is not valid")
			return
		}
		// call original func
		f(w, r)
	}
}
