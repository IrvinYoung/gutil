package aut

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"strings"
	"testing"
	"time"
)

type Account struct {
	Id       string
	Name     string
	Gender   string
	Email    string
	Password string
}

func (a *Account) DataUsingForJWT() map[string]interface{} {
	return map[string]interface{}{
		"id":   a.Id,
		"ip":   "192.168.1.100",
		"role": "guest",
		//others
	}
}

var a = &Account{
	Id:       "adsf",
	Name:     "alice",
	Gender:   "girl",
	Email:    "alice@email.com",
	Password: "12345678",
}

func TestJWT(t *testing.T) {
	testDefault(t)
	testCustom(t)
	testInvalidate(t)
	t.Log("done")
}

func testDefault(t *testing.T) {
	t.Log("-------------------------------------------default--------------------------------------------------")
	//modify expire
	MinJwtExpiration = time.Second * 5 //5 second
	if err := InitDefaultJwtEngine("1234+abcd", 6*time.Second, nil, nil); err != nil {
		t.Fatal(err)
	}

	//create a token
	JwtStr, err := NewJwt(a.DataUsingForJWT())
	if err != nil {
		t.Fatal(err)
	}
	t.Log("JWT =", JwtStr)
	//verify
	m, err := VerifyJwt(JwtStr)
	if err != nil {
		t.Log("verify failed =", err)
	} else {
		t.Logf("PASS: token content = %+v\n", m)
	}
	//wrong JWT
	m, err = VerifyJwt(strings.Replace(JwtStr, "1", "a", 1))
	if err != nil {
		t.Log("ERROR: verify failed =", err)
	}
	//expire
	time.Sleep(7 * time.Second)
	m, err = VerifyJwt(JwtStr)
	if err != nil {
		t.Log("EXPIRE: verify failed =", err)
	}
}

func testCustom(t *testing.T) {
	t.Log("-------------------------------------------custom--------------------------------------------------")
	//modify expire
	MinJwtExpiration = time.Second * 5 //5 second
	te, err := InitJwtEngine("1234+abcd", 6*time.Second, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	//create a token
	JwtStr, err := te.New(a.DataUsingForJWT())
	if err != nil {
		t.Fatal(err)
	}
	t.Log("JWT =", JwtStr)
	//verify
	m, err := te.Verify(JwtStr)
	if err != nil {
		t.Log("verify failed =", err)
	} else {
		t.Logf("PASS: token content = %+v\n", m)
	}
	//wrong JWT
	m, err = te.Verify(strings.Replace(JwtStr, "1", "a", 1))
	if err != nil {
		t.Log("ERROR: verify failed =", err)
	}
	//expire
	time.Sleep(7 * time.Second)
	m, err = te.Verify(JwtStr)
	if err != nil {
		t.Log("EXPIRE: verify failed =", err)
	}
}

func testInvalidate(t *testing.T) {
	t.Log("---------------------------------------------invalidate------------------------------------------------")
	//modify expire
	MinJwtExpiration = time.Second * 5 //5 second
	if err := InitDefaultJwtEngine("1234+abcd", 6*time.Second, Invalidate, IsScrap); err != nil {
		t.Fatal(err)
	}
	//create a token
	JwtStr, err := NewJwt(a.DataUsingForJWT())
	if err != nil {
		t.Fatal(err)
	}
	t.Log("JWT =", JwtStr)
	//verify
	m, err := VerifyJwt(JwtStr)
	if err != nil {
		t.Log("verify failed =", err)
	} else {
		t.Logf("PASS: token content = %+v\n", m)
	}
	//invalidate
	InvalidateJwt(a.Id, m["expire"])
	//expire
	m, err = VerifyJwt(JwtStr)
	if err != nil {
		t.Log("EXPIRE: verify failed =", err)
	}
}

// using redis

func Invalidate(fields ...interface{}) {
	id := fields[0].(string)
	expire := int64(fields[1].(float64))

	tm := expire - time.Now().Unix()
	if tm <= 0 {
		return
	}
	cli, err := initRedisCache()
	if err != nil {
		log.Println(err)
	}
	err = cli.Set(fmt.Sprintf("%s_%d", id, expire), "1", time.Duration(tm)*time.Second).Err()
	fmt.Println("invalidate result =",err)
}

func IsScrap(m map[string]interface{}) (err error) {
	fmt.Printf("%+v\n",m)
	id := m["id"].(string)
	expire := int64(m["expire"].(float64))
	key := fmt.Sprintf("%s_%d", id, expire)

	cli, err := initRedisCache()
	if err != nil {
		log.Println(err)
	}
	val := cli.Exists(key).Val()
	if val >= 1 {
		err = errors.New("scrap token")
		return
	}
	return
}

func initRedisCache() (cli *redis.Client, err error) {
	cli = redis.NewClient(&redis.Options{
		Addr:       "127.0.0.1:6379",
		Password:   "",
		DB:         0,
		MaxRetries: 3,
	})
	pong, err := cli.Ping().Result()
	if err != nil {
		return
	}
	if pong != "PONG" {
		err = errors.New("redis don't have a right pong result =" + pong)
		return
	}
	return
}
