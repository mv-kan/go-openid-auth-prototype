package handler

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/mv-kan/go-openid-auth-prototype/internal/utils"
	"github.com/mv-kan/go-openid-auth-prototype/openid-provider/internal"
)

var (
	loginTmpl, _ = template.New("login").Parse(`
	<!DOCTYPE html>
	<html>
		<head>
			<meta charset="UTF-8">
			<title>Login</title>
		</head>
		<body style="display: flex; align-items: center; justify-content: center; height: 100vh;">
			<form method="POST" action="/login" style="height: 200px; width: 200px;">

				<input type="hidden" name="id" value="{{.ID}}">

				<div>
					<label for="username">Username:</label>
					<input id="username" name="username" style="width: 100%">
				</div>

				<div>
					<label for="password">Password:</label>
					<input id="password" name="password" style="width: 100%">
				</div>

				<p style="color:red; min-height: 1rem;">{{.Error}}</p>

				<button type="submit">Login</button>
			</form>
		</body>
	</html>`)
)

func Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		loginGet(w, r)
		return
	case http.MethodPost:
		checkLoginPost(w, r)
		return
	default:
		utils.AllowedMethods(w, []string{http.MethodGet, http.MethodPost})
		return
	}
}

func loginGet(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	authReqIDs, ok := params["authRequestID"]
	if !ok {
		return
	}
	authReqID := authReqIDs[0]

	if !utils.ContainsID(internal.RequestStorage, authReqID) {
		return
	}
	renderLogin(w, authReqID, nil)
}

func checkLoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot parse form:%s", err), http.StatusInternalServerError)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	id := r.FormValue("id")
	err = internal.CheckUsernamePassword(username, password, id)
	if err != nil {
		renderLogin(w, id, err)
		return
	}
	authReq, err := utils.GetByID(internal.RequestStorage, id)
	if err != nil {
		renderLogin(w, id, err)
		return
	}
	url, err := authReq.GetCallbackURLAuto()
	if err != nil {
		renderLogin(w, id, err)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

func renderLogin(w http.ResponseWriter, id string, err error) {
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	data := &struct {
		ID    string
		Error string
	}{
		ID:    id,
		Error: errMsg,
	}
	err = loginTmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}