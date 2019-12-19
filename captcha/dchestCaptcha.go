package captcha

import (
	"bytes"
	"errors"
	dchestCaptcha "github.com/dchest/captcha"
	"github.com/go-redis/redis"
	"strings"
	"time"
)

type CaptchaDchest struct {
	Length int
	Width  int
	Height int
	Lang   string
	Store  dchestCaptcha.Store
}

func (c *CaptchaDchest) InitCaptcha(params ...interface{}) (instance Captcha, err error) {
	// length int
	// width int
	// height int
	// lang string	 //en,ru,zh,ja
	// store dchestCaptcha.Store

	if params == nil || len(params) != 5 {
		err = errors.New("params error")
		return
	}

	c.Length = params[0].(int)
	c.Width = params[1].(int)
	c.Height = params[2].(int)
	c.Lang = params[3].(string)
	c.Store = params[4].(dchestCaptcha.Store)

	dchestCaptcha.SetCustomStore(c.Store)
	instance = c
	return
}

func (c *CaptchaDchest) NewCaptcha(params ...interface{}) (id string, data []byte) {
	// ext string
	ext := params[0].(string)

	id = dchestCaptcha.NewLen(c.Length)
	buf := bytes.NewBuffer(nil)
	switch strings.ToLower(ext) {
	case ".png":
		dchestCaptcha.WriteImage(buf, id, c.Width, c.Height)
	case ".wav":
		dchestCaptcha.WriteAudio(buf, id, c.Lang)
	}
	data = buf.Bytes()
	return
}

func (c *CaptchaDchest) VerifyCaptcha(id, digits string) bool {
	return dchestCaptcha.VerifyString(id, digits)
}

type CaptchaDchestStore struct {
	RedisCli   *redis.Client //using redis
	Expiration time.Duration
}

// Set sets the digits for the captcha id.
func (cs *CaptchaDchestStore) Set(id string, digits []byte) {
	if id == "" || digits == nil {
		return
	}
	if cs.RedisCli == nil {
		return
	}
	cs.RedisCli.Set(id, digits, cs.Expiration)
	//log.Println("captcha set:", id, digits)
	return
}

// Get returns stored digits for the captcha id. Clear indicates
// whether the captcha must be deleted from the store.
func (cs *CaptchaDchestStore) Get(id string, clear bool) (digits []byte) {
	if cs.RedisCli == nil {
		return
	}
	digits, _ = cs.RedisCli.Get(id).Bytes()
	if clear {
		cs.RedisCli.Del(id)
	}
	//log.Println("captcha get:", id, digits, clear)
	return
}
