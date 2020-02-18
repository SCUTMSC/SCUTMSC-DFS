package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	db "../model"
	"../util"
)

const (
	pwdSalt = "*#811" // Custom number to help encode password
)

// UserSignUpGetHandler is to handle directing to signup.html
func UserSignUpGetHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Redirect(w, r, "/static/view/signup.html", http.StatusFound)
}

// UserSignUpPostHandler is to handle user signing up
func UserSignUpPostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parse the http request
	r.ParseForm()
	nickname := r.Form.Get("nickname")
	password := r.Form.Get("password")

	// Verify nickname and password
	if len(nickname) < 3 || len(password) < 6 {
		w.Write([]byte("INVALID PARAMS"))
		return
	}

	// Encode password
	password = util.Sha1([]byte(password + pwdSalt))

	// Return the http response
	if ok := db.UserSignUp(nickname, password); ok {
		w.Write([]byte("SUCCESS"))
	} else {
		w.Write([]byte("FAILED"))
	}
}

// UserSignInGetHandler is to handle directing to signup.html
func UserSignInGetHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
}

// UserSignInPostHandler is to handle user signing in
func UserSignInPostHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Parse the http request
	r.ParseForm()
	nickname := r.Form.Get("nickname")
	password := r.Form.Get("password")

	// Encode password
	password = util.Sha1([]byte(password + pwdSalt))

	// Verify nickname and password
	ok := db.UserSignIn(nickname, password)

	// Prepare response message
	resp := util.RespMsg{
		Code: 0,
		Msg:  "SUCCESS",
		Data: struct {
			Location string
			Nickname string
		}{
			Location: "/static/view/home.html",
			Nickname: nickname,
		},
	}

	// Return the http response
	if ok {
		w.Write(resp.JSONBytes())
	} else {
		w.Write([]byte("FAILED"))
	}
}

// GetUserInfoHandler is to handle sending user info to the client
func GetUserInfoHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Parse the http request
	r.ParseForm()
	nickname := r.Form.Get("nickname")

	// Retrieve user info from database
	userRecord, ok := db.GetUserInfo(nickname)

	// Prepare response message
	resp := util.RespMsg{
		Code: 0,
		Msg:  "SUCCESS",
		Data: *userRecord,
	}

	// Return the http response
	if ok {
		w.Write(resp.JSONBytes())
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}
