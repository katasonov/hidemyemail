package main

import (
	. "github.com/dchest/uniuri"
	"log"
	"net/http"
	"strings"
	//"fmt"
	"regexp"
)

func main() {

	err := LoadConfig()
	if err != nil {
		log.Fatal("Could not load config file hidemyemail.cfg")
		return
	}

	g_conn_string = g_config.DbConnectionString

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.HandleFunc("/add",
		func(w http.ResponseWriter, r *http.Request) {
			handleAdd(w, r)
		})
	http.HandleFunc("/get",
		func(w http.ResponseWriter, r *http.Request) {
			handleGet(w, r)
		})
	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			handleGetCaptcha(w, r)
		})
	log.Fatal(http.ListenAndServe(":"+g_config.Port, nil))
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	email :=  r.FormValue("email")
	matched, err := regexp.MatchString("[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?", email)
	if !matched {
		WriteIndexPage(w, true, false, email)
		return
	}
	ok, err := CheckCaptcha(r)
	if err != nil {
		WriteIndexPage(w, false, true, email)
		return
	}
	if !ok {
		WriteIndexPage(w, false, true, email)
		return
	}
	uid, err := getUidByEmailFromDatabase(email)
	if err == nil {
		WriteSecureLinkPage(w, uid)
		return
	}
	uid = NewLen(8)
	err = addEmailToDatabase(uid, email)
	if err != nil {
		w.Write([]byte("Db Error occurred: " + err.Error()))
		return
	}

	WriteSecureLinkPage(w, uid)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	uid := r.FormValue("email_uid")
	ok, err := CheckCaptcha(r)
	if err != nil {
		WriteAccessEmailPage(w, uid, true)
		return
	}
	if !ok {
		WriteAccessEmailPage(w, uid, true)
		return
	}
	email, err := getEmailByUidFromDatabase(uid)
	if err != nil {
		WriteEmailNotFoundPage(w, uid)
		return
	}
	WriteEmailPage(w, email)
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
