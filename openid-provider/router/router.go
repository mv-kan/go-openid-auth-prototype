package router

import (
	"net/http"

	"github.com/mv-kan/go-openid-auth-prototype/internal/vars"
	"github.com/mv-kan/go-openid-auth-prototype/openid-provider/handler"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(vars.AUTHN_ENDPOINT, handler.Authenticate)
	mux.HandleFunc(vars.LOGIN_ENDPOINT, handler.Login)
	mux.HandleFunc(vars.TOKEN_ENDPOINT, handler.Token)
	mux.HandleFunc(vars.CHECK_TOKEN_ENDPOINT, handler.CheckToken)
	return mux
}
