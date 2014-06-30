package main

import (
	. "github.com/dchest/uniuri"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	//"fmt"
	"bytes"
	. "github.com/katasonov/asycache"
)

//var templates = template.Must(template.ParseGlob("html/*"))
var g_index_templ = template.New("base")
var g_captcha_templ = template.New("base")
var g_url_templ = template.New("base")
var g_email_templ = template.New("base")

//var g_index_templ = template.New("index").ParseFiles("html/base.html")

func main() {
	g_index_templ, _ = g_index_templ.ParseFiles("html/index.html")
	g_index_templ, _ = g_index_templ.ParseFiles("html/base.html")

	g_captcha_templ, _ = g_captcha_templ.ParseFiles("html/captcha.html")
	g_captcha_templ, _ = g_captcha_templ.ParseFiles("html/base.html")

	g_url_templ, _ = g_url_templ.ParseFiles("html/url.html")
	g_url_templ, _ = g_url_templ.ParseFiles("html/base.html")

	g_email_templ, _ = g_email_templ.ParseFiles("html/email.html")
	g_email_templ, _ = g_email_templ.ParseFiles("html/base.html")

	g_conn_string = "hidemyemail:Avk241083@/hidemyemaildb"

	c := MakeCache(10 * time.Minute)

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.HandleFunc("/add",
		func(w http.ResponseWriter, r *http.Request) {
			handleAdd(w, r, c)
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

func handleAdd(w http.ResponseWriter, r *http.Request, c Cache) {
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

func parseTemplate(file string, data interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	//t := template.New(file).Funcs(funcMaps)
	t := template.New(file)
	baseBytes, err := ioutil.ReadFile("html/base.html")

	if err != nil {
		return nil, err
	}
	t, err = t.Parse(string(baseBytes))
	if err != nil {
		return nil, err
	}

	t, err = t.ParseFiles("html/" + file)
	if err != nil {
		return nil, err
	}
	err = t.Execute(buf, data)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
