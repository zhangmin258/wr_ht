package utils

import (
	"net/url"
	"strconv"
	"zcm_tools/googauth"
)

func Authenticate(uid int, verify_code string) (bool, error) {
	secret := googauth.CreateSecret2(LoginVerifyCodePrefox, strconv.Itoa(uid))
	otpconf := &googauth.OTPConfig{
		Secret:     secret,
		WindowSize: 3,
	}
	return otpconf.Authenticate(verify_code)
}

func CreateXjdSecret(uid int) string {
	return googauth.CreateSecret2(LoginVerifyCodePrefox, strconv.Itoa(uid))
}

func CreateXjdAuthURLEscape(secret, username string) string {
	company := "微融"
	if RunMode == "release" {
	} else if RunMode == "test" {
		company += "测试服"
	} else {
		company += "开发服"
	}
	return url.QueryEscape(googauth.CreateAuthURL(secret, company, username))
}
