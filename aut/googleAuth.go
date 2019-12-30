package aut

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"
)

func CreateGoogleAuthSecret() (secret string, err error) {
	rand.Seed(time.Now().UnixNano())
	part1, part2 := rand.Int63(), rand.Int63()

	b := bytes.NewBuffer(nil)
	if err = binary.Write(b, binary.LittleEndian, part1); err != nil {
		return
	}
	b1 := b.Bytes()[0:5]

	b.Reset()
	if err = binary.Write(b, binary.BigEndian, part2); err != nil {
		return
	}
	b2 := b.Bytes()[0:5]
	key := append(b1, b2...)
	//log.Println("key =", key)
	secret = base32.StdEncoding.EncodeToString(key)
	return
}

func makeGoogleAuthCode(secret string, t int64) (code uint32, err error) {
	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return
	}

	hash := hmac.New(sha1.New, key)
	err = binary.Write(hash, binary.BigEndian, t)
	if err != nil {
		return
	}

	h := hash.Sum(nil)
	offset := h[len(h)-1] & 0x0f
	code = binary.BigEndian.Uint32(h[offset : offset+4])
	code &= 0x7fffffff
	code %= 1000000
	return
}

//VerifyGoogleAuthByTime TOTP
func VerifyGoogleAuthByTime(secret string, code uint32) (err error) {
	target, err := makeGoogleAuthCode(secret, time.Now().Unix()/30)
	if err != nil {
		return
	}
	if target != code {
		err = errors.New("invalid code")
		return
	}
	return
}

//VerifyGoogleAuthByCounter HOTP
func VerifyGoogleAuthByCounter(secret string, code uint32, counter int64) (err error) {
	target, err := makeGoogleAuthCode(secret, counter+1)
	if err != nil {
		return
	}
	if target != code {
		err = errors.New("invalid code")
		return
	}
	return
}

func MakeGoogleAuthURI(secret, user, issuer string, counter int64) (uri string, err error) {
	var (
		verifyType string
		params     = make(url.Values)
	)
	if counter <= 0 {
		verifyType = "totp"
	} else {
		verifyType = "hotp"
		params.Set("counter", strconv.FormatInt(counter, 10))
	}
	params.Set("secret", secret)
	if issuer != "" {
		params.Set("issuer", issuer)
	}

	uri = fmt.Sprintf("otpauth://%s/%s:%s?%s", verifyType, issuer, user, params.Encode())
	return
}
