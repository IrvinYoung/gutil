package captcha

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestDchestCaptcha(t *testing.T) {
	cli, err := initRedisCache()
	if err != nil {
		t.Fatal(err)
	}

	cs := &CaptchaDchestStore{
		RedisCli:   cli,
		Expiration: time.Hour,
	}

	c := &CaptchaDchest{}
	ct, err := c.InitCaptcha(5, 240, 80, "en", cs)
	if err != nil {
		t.Fatal(err)
	}

	//make
	id, data := ct.NewCaptcha(".png", "en")
	t.Logf("id=%s\tdata=./%s.png\n", id, id)

	//show
	err = ioutil.WriteFile("./"+id+".png", data, os.FileMode(0644))
	if err != nil {
		t.Fatal(err)
	}

	//check
	var digits string
	fmt.Scan(&digits)	//FIXME:

	if ct.VerifyCaptcha(id, digits) {
		t.Log("PASS:", id, digits)
	} else {
		t.Log("FAIL:", id, digits)
	}
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
