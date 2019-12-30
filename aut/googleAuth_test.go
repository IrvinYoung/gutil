package aut

import (
	"testing"
	"time"
)

func TestCreateSecretAndCode(t *testing.T) {
	secret, err := CreateGoogleAuthSecret()
	t.Log(secret, err)

	code, err := makeGoogleAuthCode(secret, time.Now().Unix())
	t.Logf("1 -> %06d, %v\n", code, err)
	time.Sleep(time.Second * 2)
	code, err = makeGoogleAuthCode(secret, time.Now().Unix())
	t.Logf("2 -> %06d, %v\n", code, err)
}

func TestTOTP(t *testing.T) {
	secret := "JTVCT5XKJTVCT5XK"
	code, _ := makeGoogleAuthCode(secret, 1/30)
	t.Log("code(second 1) =", code, code == 695775) //695775

	code, _ = makeGoogleAuthCode(secret, 10000/30)
	t.Log("code(second 10000) =", code, code == 323664) //323664
}

func TestHOTP(t *testing.T) {
	secret := "JTVCT5XKJTVCT5XK"
	code, _ := makeGoogleAuthCode(secret, 2)
	t.Log("code(second 1) =", code) //479349
	if err := VerifyGoogleAuthByCounter(secret, code, 1); err != nil {
		t.Log("FAIL", err)
	} else {
		t.Log("PASS")
	}

	code, _ = makeGoogleAuthCode(secret, 10001)
	t.Log("code(second 10000) =", code)
	if err := VerifyGoogleAuthByCounter(secret, code, 10000); err != nil {
		t.Log("FAIL", err)
	} else {
		t.Log("PASS")
	}
}

func TestMakeURI(t *testing.T){
	totpURI,err := MakeGoogleAuthURI("JTVCT5XKJTVCT5XK", "Admin", "GUtil", 0)
	t.Log(totpURI,err)

	hotpURI,err := MakeGoogleAuthURI("JTVCT5XKJTVCT5XK", "Admin", "GUtil", 1)
	t.Log(hotpURI,err)
}