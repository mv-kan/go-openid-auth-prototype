package handler

import (
	"net/http"

	"github.com/mv-kan/go-openid-auth-prototype/internal/log"
)

func GetProtectedSuperSecret(w http.ResponseWriter, r *http.Request) {
	log.Info("get protected super secret is called")
	w.Write([]byte("Hi, it is sure really cool access token you have here"))
}
