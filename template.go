package main

import (
	"net/http"
	"html/template"
	"github.com/dchest/captcha"
)

type IndexPageValues struct {
	EmailError bool
	CaptchaError bool
	Email string //contains email that should be inserted in case of error.
}

type AccessEmailPageValues struct {
	Key string
	CaptchaError bool
	CaptchaId string
}

type EmailPageValues struct {
	Email string
	IsEmail bool
	EmailRef string
}

func WriteIndexPage(w http.ResponseWriter, email_error bool, captcha_error bool, email string) {
	writeHtmlWithValues(w, "index.html", &IndexPageValues{email_error, captcha_error, email})
}

func WriteAccessEmailPage(w http.ResponseWriter, uid string, captcha_error bool) {
	writeHtmlWithValues(w, "captcha.html", &AccessEmailPageValues{uid, captcha_error, captcha.New()})
}

func WriteEmailPage(w http.ResponseWriter, email string) {
	writeHtmlWithValues(w, "email.html", &EmailPageValues{Email: email, IsEmail: true})
}

func WriteUrlPage(w http.ResponseWriter, email string) {
	ref := email
	if !urlHasPrefix(email) {
		ref = "http://" + ref
	}
	writeHtmlWithValues(w, "email.html", &EmailPageValues{Email: email, IsEmail: false, EmailRef: ref})
}

func WriteSecureLinkPage(w http.ResponseWriter, uid string) {
	host := "http://" + g_config.Host
	if g_config.Port != "80" && g_config.ShowPort != 0 {
		host = host + ":" + g_config.Port
	}
	writeHtmlWithValues(w, "url.html", &struct{ Host string; Key string }{host, uid})
}

func WriteEmailNotFoundPage(w http.ResponseWriter, uid string) {
	writeHtmlWithValues(w, "emailnotfound.html", &struct{Key string}{uid})
}

func writeHtmlWithValues(w http.ResponseWriter, file string, data interface{}) {
	tmpl := template.New("base")
	var err error
	if tmpl, err = tmpl.ParseFiles(g_config.ResourcePath + "/html/base.html"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if tmpl, err = tmpl.ParseFiles(g_config.ResourcePath + "/html/" + file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	type Base struct {
		Content interface{}
	}
	if err = tmpl.Execute(w, &Base{Content: data}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
