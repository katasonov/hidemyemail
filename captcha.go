package main

import (
	"net/http"
	"github.com/dchest/captcha"
)

func CheckCaptcha(r *http.Request) (bool) {
	return captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution"))
}
