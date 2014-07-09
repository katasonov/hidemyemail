package main

import (
	"regexp"
)

func isEmail(v string)bool {
	matched, _ := regexp.MatchString("[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?", v)
	if !matched {
		return false
	}
	return true
}

func isUrl(v string)bool{
	matched, _ := regexp.MatchString("^((https?|ftp|file)://*.)", v)
	return matched
}
