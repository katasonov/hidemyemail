package main

import (
	"net/http"
	"net/url"
	"io/ioutil"
)

func CheckCaptcha(r *http.Request) (bool, error) {
	recaptcha_challenge_field := r.FormValue("recaptcha_challenge_field")
	recaptcha_response_field := r.FormValue("recaptcha_response_field")
	resp, err := http.PostForm(
		"http://www.google.com/recaptcha/api/verify",
		url.Values{
		"privatekey": {"6LdUtvUSAAAAAEuQtd3u6vaSeXVNZV2k9A1R_XG7"},
		"remoteip":   {r.RemoteAddr},
		"challenge":  {recaptcha_challenge_field},
		"response":   {recaptcha_response_field}})
	if err != nil {
		return false, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	resultS := string(body[:4])
	if resultS != "true" {
		return false, nil
	}
	return true, nil
}
