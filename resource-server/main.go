package main

import (
	"log"
	"net/http"

	"github.com/mv-kan/go-openid-auth-prototype/resource-server/middleware"
)

var (
	host = "localhost:7000"
)

func main() {
	http.HandleFunc("/protected", middleware.OnlyMethod(http.MethodGet, []string{http.MethodGet}, middleware.ValidateToken(middleware.GetProtectedSuperSecret)))
	log.Println("Start serving on " + host + " ...")
	log.Fatal(http.ListenAndServe(host, nil))
}
