package main

import (
	"net/http"
	"strings"
	. "github.com/dchest/uniuri"
)

func handleAdd(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	email :=  r.FormValue("email")
	if len(email) > 127 {
		WriteIndexPageWithInvalidEmailLen(w, email)
		return
	}
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)
	uid := NewLen(8)
	err := addEmailToDatabase(uid, email)
	if err != nil {
		w.Write([]byte("Db Error occurred: " + err.Error()))
		return
	}

	WriteSecureLinkPage(w, uid)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	uid := r.FormValue("email_uid")
	ok := CheckCaptcha(r)
	if !ok {
		WriteAccessEmailPage(w, uid, true)
		return
	}
	email, err := getEmailByUidFromDatabase(uid)
	if err != nil {
		WriteEmailNotFoundPage(w, uid)
		return
	}

	if !isEmail(email) {
		WriteUrlPage(w, email)
	} else {
		WriteEmailPage(w, email)
	}
}

func handleGetCaptcha(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	key := strings.TrimLeft(r.URL.Path, "/")
	if key == "" {
		//index page
		WriteIndexPage(w, false, false, "")
		return
	}
	WriteAccessEmailPage(w, key, false)
}

