package handler

import "net/http"

func GetProtectedSuperSecret(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi, it is sure really cool access token you have here"))
}
