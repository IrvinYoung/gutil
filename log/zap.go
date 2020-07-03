package log

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"net/url"
	"strconv"
	"strings"
)

type ZapSkinRedisQueue struct {
	cli      *redis.Client
	queue    string
	fromLeft bool
}

func (s ZapSkinRedisQueue) Sync() error {
	//s.cli.Sync()	//not implemented
	return nil
}

func (s ZapSkinRedisQueue) Close() error {
	if s.cli == nil {
		return nil
	}
	return s.cli.Close()
}

func (s ZapSkinRedisQueue) Write(p []byte) (n int, err error) {
	if s.fromLeft {
		err = s.cli.LPush(s.queue, string(p)).Err()
	} else {
		err = s.cli.RPush(s.queue, string(p)).Err()
	}
	return len(p), err
}

func NewZapSkinRedisQueue(url *url.URL) (sink zap.Sink, err error) {
	if url.Host == "" {
		err = errors.New("lost redis host")
		return
	}

	s := ZapSkinRedisQueue{}
	query := url.Query()
	s.queue = strings.ToLower(query.Get("queue"))
	if s.queue == "" {
		err = errors.New("lost redis queue name")
		return
	}
	operate := strings.ToLower(query.Get("op"))
	if operate != "lpush" && operate != "rpush" {

	}
	switch operate {
	case "lpush":
		s.fromLeft = true
	case "rpush":
		s.fromLeft = false
	default:
		err = fmt.Errorf("unsupported redis operate: %s", operate)
		return
	}
	//init redis client
	db, err := strconv.Atoi(query.Get("db"))
	if err != nil {
		err = fmt.Errorf("invalid redis db number: %s", query.Get("db"))
		return
	}
	s.cli = redis.NewClient(&redis.Options{
		Addr:     url.Host,
		Password: query.Get("password"),
		DB:       db,
	})
	if _, err = s.cli.Ping().Result(); err != nil {
		return
	}
	if err == redis.Nil {
		err = errors.New("redis error")
	}
	return s, err
}

type ZapSkinRabbitMQ struct {
}

type ZapSkinMongo struct {
}

type ZapSkinXorm struct {
}
