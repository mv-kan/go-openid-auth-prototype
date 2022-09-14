package middleware

import (
	"net/http"
	"strings"

	"github.com/mv-kan/go-openid-auth-prototype/internal/utils"
)

func contains[T comparable](sl []T, elem T) bool {
	for _, value := range sl {
		if value == elem {
			return true
		}
	}
	return false
}

func AllowedMethods(allowedMethods []string, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if contains(allowedMethods, r.Method) {
			f(w, r)
		} else {
			w.Header().Add("Allow", strings.Join(allowedMethods, ", "))
			utils.WriteResponse(w, http.StatusMethodNotAllowed, "Not allowed")
		}
	}
}
