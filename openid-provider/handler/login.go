package handler

import (
	"errors"
	"fmt"
	"net/http"
	"text/template"

	"github.com/mv-kan/go-openid-auth-prototype/internal/log"
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

// render login page
func loginGet(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	authReqIDs, ok := params["authRequestID"]
	if !ok {
		log.Debug("not parsed properly")
		utils.WriteResponse(w, http.StatusInternalServerError, "")
		return
	}
	authReqID := authReqIDs[0]

	_, err := utils.GetByID(internal.RequestStorage, authReqID)
	if errors.Is(err, utils.ErrNotFound) {
		log.Debug("authenticate request does not exist")
		utils.ResponseJSON(w, http.StatusForbidden, map[string]string{"error": "auth request does not exist"})
		return
	} else if err != nil {
		log.Error(err.Error())
		utils.WriteResponse(w, http.StatusInternalServerError, "")
		return
	}
	renderLogin(w, authReqID, nil)
}

// check login post info
func checkLoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, fmt.Sprintf("cannot parse form:%s", err), http.StatusInternalServerError)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	id := r.FormValue("id")
	err = internal.CheckUsernamePassword(username, password, id)
	if err != nil {
		log.Error(err.Error())
		renderLogin(w, id, err)
		return
	}
	authReq, err := utils.GetByID(internal.RequestStorage, id)
	if err != nil {
		log.Debug(err.Error())
		renderLogin(w, id, err)
		return
	}
	url, err := authReq.GetCallbackURLAuto()
	if err != nil {
		// We use error function because this is out of reach to user, it purely server side thing
		log.Error(err.Error())
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
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
