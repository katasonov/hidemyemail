package main

import (
	. "github.com/dchest/uniuri"
	"html/template"
	"io/ioutil"
	"log"
	. "net/http"
	"net/url"
	"strings"
	"time"
	//"fmt"
	. "github.com/katasonov/asycache"
)

func main() {

	c := MakeCache(10 * time.Minute)

	HandleFunc("/add",
		func(w ResponseWriter, r *Request) {
			handleAdd(w, r, c)
		})
	HandleFunc("/get",
		func(w ResponseWriter, r *Request) {
			handleGet(w, r)
		})
	HandleFunc("/",
		func(w ResponseWriter, r *Request) {
			handleGetCaptcha(w, r)
		})
	log.Fatal(ListenAndServe(":8080", nil))
}

func handleAdd(w ResponseWriter, r *Request, c Cache) {
	w.WriteHeader(StatusOK)

	values := r.URL.Query()
	err := addEmailToDatabase(NewLen(8), values.Get("email"))
	if err != nil {
		w.Write([]byte("Db Error occurred: " + err.Error()))
		return
	}
	w.Write([]byte("OK"))
}

func handleGet(w ResponseWriter, r *Request) {
	w.WriteHeader(StatusOK)
	recaptcha_challenge_field := r.FormValue("recaptcha_challenge_field")
	recaptcha_response_field := r.FormValue("recaptcha_response_field")
	uid := r.FormValue("email_uid")
	resp, err := PostForm(
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

func handleGetCaptcha(w ResponseWriter, r *Request) {
	w.WriteHeader(StatusOK)
	key := strings.TrimLeft(r.URL.Path, "/")
	if key == "" {
		//index page
		writeHtmlWithValues(w, "index.html", &struct{}{})
		return
	}
	writeHtmlWithValues(w, "captcha.html", &struct{ Key string }{key})
}

func writeHtmlWithValues(w ResponseWriter, html_file string, data interface{}) {
	t, err := template.ParseFiles("html/" + html_file)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		t.Execute(w, &data)
	}
}

/*
func getEmailFromDatabase(key string) (*CacheEntity, bool){
	con, err := sql.Open("mysql56", "hidemyemails:Avk241083@/hidemyemailsdb")
	defer con.Close()
	row := con.QueryRow("select data as email from email where uid=?", key)
	entity := &CacheEntity{uid: key}
	err = row.Scan(&entity.email)
	if err != nil {
		return nil, false
	}
	return entity, true
}
*/
