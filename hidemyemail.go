package main

import (
	. "github.com/dchest/uniuri"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	//"fmt"
)

func main() {

	g_conn_string = "hidemyemail:Avk241083@/hidemyemaildb"

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
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	values := r.URL.Query()
	email := values.Get("email")
	uid, err := getUidByEmailFromDatabase(email)
	if err == nil {
		writeHtmlWithValues(w, "url.html", &struct{ Key string }{uid})
		return
	}
	uid = NewLen(8)
	err = addEmailToDatabase(uid, email)
	if err != nil {
		w.Write([]byte("Db Error occurred: " + err.Error()))
		return
	}
	writeHtmlWithValues(w, "url.html", &struct{ Key string }{uid})
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	recaptcha_challenge_field := r.FormValue("recaptcha_challenge_field")
	recaptcha_response_field := r.FormValue("recaptcha_response_field")
	uid := r.FormValue("email_uid")
	resp, err := http.PostForm(
		"http://www.google.com/recaptcha/api/verify",
		url.Values{
			"privatekey": {"6LdUtvUSAAAAAEuQtd3u6vaSeXVNZV2k9A1R_XG7"},
			"remoteip":   {r.RemoteAddr},
			"challenge":  {recaptcha_challenge_field},
			"response":   {recaptcha_response_field}})
	if err != nil {
		w.Write([]byte("Failed to get captcha result from google"))
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.Write([]byte("Failed to read response body"))
		return
	}
	resultS := string(body[:4])
	if resultS != "true" {
		w.Write([]byte("Invalid captcha"))
		return
	}
	email, err := getEmailByUidFromDatabase(uid)
	if err != nil {
		w.Write([]byte("Db Error occurred: " + err.Error()))
		return
	}
	writeHtmlWithValues(w, "email.html", &struct{ Email string }{email})
}

func handleGetCaptcha(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	key := strings.TrimLeft(r.URL.Path, "/")
	if key == "" {
		//index page
		writeHtmlWithValues(w, "index.html", &struct{}{})
		return
	}
	writeHtmlWithValues(w, "captcha.html", &struct{ Key string }{key})
}

func writeHtmlWithValues(w http.ResponseWriter, file string, data interface{}) {
	//err := templates.ExecuteTemplate(w, html_file, &data)
	//g_index_templ.Parse(w, "html/index.html")
	//err := t.Execute(w, data)
	//bval, err := parseTemplate(file, data)
	tmpl := template.New("base")
	var err error
	if tmpl, err = tmpl.ParseFiles("html/base.html"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if tmpl, err = tmpl.ParseFiles("html/" + file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	type Base struct {
		Content interface{}
	}
	if err = tmpl.Execute(w, &Base{Content: data}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
