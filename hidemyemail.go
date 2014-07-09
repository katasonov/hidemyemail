package main

import (
	"log"
	"net/http"
	"github.com/dchest/captcha"
)

var chttp = http.NewServeMux()

func main() {

	err := LoadConfig()
	if err != nil {
		log.Fatal("Could not load config file hidemyemail.cfg")
		return
	}

	g_conn_string = g_config.DbConnectionString

	http.Handle("/captcha/", captcha.Server(captcha.StdWidth, captcha.StdHeight))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(g_config.ResourcePath + "/images"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(g_config.ResourcePath + "/css"))))
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
