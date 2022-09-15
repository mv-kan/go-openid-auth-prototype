package router

import (
	"net/http"

	"github.com/mv-kan/go-openid-auth-prototype/internal/vars"
	"github.com/mv-kan/go-openid-auth-prototype/openid-provider/handler"
)

func New() {
	mux := http.NewServeMux()
	mux.HandleFunc(vars.AUTHN_ENDPOINT, handler.Authenticate)
	mux.HandleFunc(vars.LOGIN_ENDPOINT, handler.Login)
}
