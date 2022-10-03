package main

import (
	"net/http"

	"github.com/mv-kan/go-openid-auth-prototype/internal/log"
	"github.com/mv-kan/go-openid-auth-prototype/internal/vars"
	"github.com/mv-kan/go-openid-auth-prototype/relying-party/handler"
)

func main() {
	http.HandleFunc("/callback", handler.Callback)
	http.HandleFunc("/", handler.Index)
	log.Info("Start relying party serving on " + vars.RP_HOST + " ...")
	log.Error(http.ListenAndServe(vars.RP_HOST, nil).Error())
}
