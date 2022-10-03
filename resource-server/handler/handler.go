package handler

import (
	"net/http"

	"github.com/mv-kan/go-openid-auth-prototype/internal/log"
)

func GetProtectedSuperSecret(w http.ResponseWriter, r *http.Request) {
	log.Info("get protected super secret is called")
	w.Write([]byte("This is very secret message from protected resource, only authorizated and authenticated can see it."))
}
