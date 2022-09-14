package main

import (
	"log"
	"net/http"

	internalMiddleware "github.com/mv-kan/go-openid-auth-prototype/internal/middleware"
	"github.com/mv-kan/go-openid-auth-prototype/internal/vars"
	"github.com/mv-kan/go-openid-auth-prototype/resource-server/handler"
	"github.com/mv-kan/go-openid-auth-prototype/resource-server/middleware"
)

func main() {
	http.HandleFunc("/protected", internalMiddleware.AllowedMethods([]string{http.MethodGet}, middleware.ValidateToken(handler.GetProtectedSuperSecret)))
	log.Println("Start serving on " + vars.RS_HOST + " ...")
	log.Fatal(http.ListenAndServe(vars.RS_HOST, nil))
}
