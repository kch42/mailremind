package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/gorilla/sessions"
	"kch42.de/gostuff/mailremind/model"
	"log"
	"net/http"
)

type loginTpldata struct {
	Error, Success string
}

type loginFormdata struct {
	Mail, Password string
}

func login(user model.User, sess *sessions.Session, req *http.Request) (interface{}, model.User) {
	outdata := &loginTpldata{}

	if user != nil {
		outdata.Success = "You are already logged in"
		return outdata, user
	}

	if req.Method != "POST" {
		return outdata, user
	}

	if err := req.ParseForm(); err != nil {
		outdata.Error = "Data of form could not be understand. If this happens again, please contact support!"
		return outdata, user
	}

	indata := new(loginFormdata)
	if err := formdec.Decode(indata, req.Form); (err != nil) || (indata.Mail == "") || (indata.Password == "") {
		outdata.Error = "Input data wrong or missing. Please fill in all values."
		return outdata, user
	}

	user, err := dbcon.UserByMail(indata.Mail)
	switch err {
	case nil:
	case model.NotFound:
		outdata.Error = "E-Mail or password was wrong."
		return outdata, user
	default:
		log.Printf("Error while loding user data (login): %s", err)
		outdata.Error = "User data could not be loaded. Please contact support, if this happens again."
		return outdata, user
	}

	if bcrypt.CompareHashAndPassword(user.PWHash(), []byte(indata.Password)) != nil {
		outdata.Error = "E-Mail or password was wrong."
		return outdata, user
	}

	sess.Values["uid"] = user.ID().String()
	outdata.Success = "Login successful"
	return outdata, user
}

func logout(user model.User, sess *sessions.Session, req *http.Request) (interface{}, model.User) {
	delete(sess.Values, "uid")
	return &msgTpldata{Class: "success", Title: "Logout", Msg: "Your are now logged out."}, nil
}
