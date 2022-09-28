package main

import (
	"fmt"
	"net/http"

	"github.com/mv-kan/go-openid-auth-prototype/internal/log"
	"github.com/mv-kan/go-openid-auth-prototype/internal/vars"
)

func main() {
	http.HandleFunc("/callback", callbackHandler)
	log.Error(http.ListenAndServe(vars.RP_HOST, nil).Error())
}
func callbackHandler(w http.ResponseWriter, r *http.Request) {
	// gather information about request and log it
	uri := r.URL.String()
	method := r.Method
	fmt.Printf("uri = %s, method = %s", uri, method)
	w.WriteHeader(http.StatusNotImplemented)
}
