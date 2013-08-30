package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/gorilla/sessions"
	"kch42.de/gostuff/mailremind/model"
	"log"
	"net/http"
)

type pwresetTpldata struct {
	Error, Success, Code, UID string
}

func pwreset(user model.User, sess *sessions.Session, req *http.Request) interface{} {
	if err := req.ParseForm(); err != nil {
		return &pwresetTpldata{Error: "Could not understand formdata."}
	}

	code := req.FormValue("Code")
	_uid := req.FormValue("U")
	pw1 := req.FormValue("Password")
	pw2 := req.FormValue("PasswordAgain")

	if code == "" {
		return &pwresetTpldata{Error: "Wrong password reset code"}
	}

	uid, err := db.ParseDBID(_uid)
	if err != nil {
		return &pwresetTpldata{Error: "Invalid user ID"}
	}

	if user, err = dbcon.UserByID(uid); err != nil {
		return &pwresetTpldata{Error: "User not found"}
	}

	if user.ActivationCode() != code {
		return &pwresetTpldata{Error: "Wrong activation code"}
	}

	outdata := &pwresetTpldata{UID: _uid, Code: code}

	if req.Method != "POST" {
		return outdata
	}

	if pw1 == "" {
		outdata.Error = "Password must not be empty."
		return outdata
	}

	if pw1 != pw2 {
		outdata.Error = "Passwords are not identical."
		return outdata
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pw1), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Could not has password: %s", err)
		outdata.Error = "Failed hashing you password. If this happens again, please contact support."
		return outdata
	}

	if err := user.SetPWHash(hash); err != nil {
		log.Printf("Error while storing new password: %s", err)
		outdata.Error = "Could not store password. If this happens again, please contact support."
		return outdata
	}

	if err := user.SetActivationCode(""); err != nil {
		log.Printf("Error resetting acCode: %s", err)
	}

	outdata.Success = "Password was changed"
	return outdata
}

type forgotpwTpldata struct {
	Error, Success string
}

func forgotpw(user model.User, sess *sessions.Session, req *http.Request) interface{} {
	if req.Method != "POST" {
		return &forgotpwTpldata{}
	}

	if err := req.ParseForm(); err != nil {
		return &forgotpwTpldata{Error: "Could not understand formdata."}
	}

	email := req.FormValue("Mail")
	if email == "" {
		return &forgotpwTpldata{Error: "E-Mail must not be empty."}
	}

	user, err := dbcon.UserByMail(email)
	if err != nil {
		return &forgotpwTpldata{Error: "E-Mail not found."}
	}

	key := genAcCode()
	if err := user.SetActivationCode(key); err != nil {
		log.Printf("Could not store pwreset key: %s", err)
		return &forgotpwTpldata{Error: "Could not store keyword reset code. If this happens again, please contact support."}
	}

	if !SendPwresetLink(user.Email(), key, user.ID()) {
		return &forgotpwTpldata{Error: "Could not send reset E-Mail. If this happens again, please contact support."}
	}

	return &forgotpwTpldata{Success: "We sent you an E-Mail with further instructions."}
}
