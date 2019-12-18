package aut

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JwtEngine struct {
	Sign       string
	Expiration time.Duration

	Invalidate func(...interface{})
	IsScrap    func(map[string]interface{}) error
}

var (
	MinJwtSignLength = 8
	MinJwtExpiration = time.Hour
)

var (
	defaultJwtEngine     *JwtEngine
	defaultJwtExpiration time.Duration //second
)

func InitDefaultJwtEngine(
	sign string, expiration time.Duration,
	fInvalidate func(...interface{}),
	fIsScrap func(map[string]interface{}) error) (err error) {
	defaultJwtEngine, err = InitJwtEngine(sign, expiration, fInvalidate, fIsScrap)
	return
}

func InitJwtEngine(
	sign string, expiration time.Duration,
	fInvalidate func(...interface{}),
	fIsScrap func(map[string]interface{}) error) (t *JwtEngine, err error) {
	if sign == "" || len(sign) < MinJwtSignLength {
		err = errors.New("jwt sign is invalid")
		return
	}
	if expiration < MinJwtExpiration {
		err = errors.New("jwt expiration need more than 60 minutes")
		return
	}
	t = &JwtEngine{
		Sign:       sign,
		Expiration: expiration,
		Invalidate: fInvalidate,
		IsScrap:    fIsScrap,
	}
	return
}

func NewJwt(m map[string]interface{}) (JwtStr string, err error) {
	if defaultJwtEngine == nil {
		err = errors.New("default jwt engine is invalid")
		return
	}
	JwtStr, err = defaultJwtEngine.New(m)
	return
}

func VerifyJwt(JwtStr string) (m map[string]interface{}, err error) {
	if defaultJwtEngine == nil {
		err = errors.New("default jwt engine is invalid")
		return
	}
	m, err = defaultJwtEngine.Verify(JwtStr)
	return
}

func (t *JwtEngine) New(m map[string]interface{}) (JwtStr string, err error) {
	mc := jwt.MapClaims(m)
	mc["expire"] = time.Now().Add(t.Expiration).Unix()
	JwtStr, err = jwt.NewWithClaims(jwt.SigningMethodHS256, mc).SignedString([]byte(t.Sign))
	return
}

func (t *JwtEngine) Verify(JwtStr string) (m map[string]interface{}, err error) {
	if JwtStr == "" {
		err = errors.New("token is nil")
		return
	}
	token, err := jwt.Parse(JwtStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(t.Sign), nil
	})
	if err != nil {
		return
	}
	//1.check internal data
	if !token.Valid {
		err = jwt.ErrInvalidKey
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("get claims from jwt error")
		return
	}
	expire, ok := claims["expire"]
	if !ok {
		err = errors.New("token lost expire")
		return
	}
	if int64(expire.(float64)) < time.Now().Unix() {
		err = errors.New("token expire")
		return
	}
	//2.check recode data
	if t.IsScrap != nil {
		if err = t.IsScrap(claims); err != nil {
			return
		}
	}
	m = claims
	return
}

func InvalidateJwt(params ...interface{}) {
	if defaultJwtEngine == nil || defaultJwtEngine.Invalidate == nil {
		return
	}
	defaultJwtEngine.Invalidate(params...)
}
