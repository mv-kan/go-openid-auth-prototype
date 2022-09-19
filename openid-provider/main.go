package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/mv-kan/go-openid-auth-prototype/internal/vars"
	"github.com/mv-kan/go-openid-auth-prototype/openid-provider/internal"
	"github.com/mv-kan/go-openid-auth-prototype/openid-provider/router"
)

func RegisterClients() {
	redirectURI, err := url.JoinPath(vars.RP_HOST, vars.REDIRECT_URI)
	if err != nil {
		log.Fatal(err)
	}
	client := internal.Client{
		ID:           "web",
		Secret:       "secret",
		RedirectURIs: []string{redirectURI},
	}
	internal.ClientStorage = append(internal.ClientStorage, client)
}

func RegisterUsers() {
	user := internal.User{
		Username: "username",
		Password: "password",
	}
	internal.UserStorage = append(internal.UserStorage, user)
}

func main() {
	RegisterUsers()
	RegisterClients()
	mux := router.New()
	log.Println("Start serving on " + vars.OP_HOST + " ...")
	log.Fatal(http.ListenAndServe(vars.OP_HOST, mux))
}
