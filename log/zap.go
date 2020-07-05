package log

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"net/url"
	"strconv"
	"strings"
)

type ZapSinkRedisQueue struct {
	cli      *redis.Client
	queue    string
	fromLeft bool
}

func (s ZapSinkRedisQueue) Sync() error {
	//s.cli.Sync()	//not implemented
	return nil
}

func (s ZapSinkRedisQueue) Close() error {
	if s.cli == nil {
		return nil
	}
	return s.cli.Close()
}

func (s ZapSinkRedisQueue) Write(p []byte) (n int, err error) {
	if s.fromLeft {
		err = s.cli.LPush(s.queue, string(p)).Err()
	} else {
		err = s.cli.RPush(s.queue, string(p)).Err()
	}
	return len(p), err
}

func NewZapSinkRedisQueue(url *url.URL) (sink zap.Sink, err error) {
	if url.Host == "" {
		err = errors.New("lost redis host")
		return
	}

	s := ZapSinkRedisQueue{}
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

type ZapSinkRabbitMQ struct {
	con       *amqp.Connection
	ch        *amqp.Channel
	errorChan chan *amqp.Error

	levelKey     string
	exchangeName string
}

func (s ZapSinkRabbitMQ) Sync() error {
	return nil
}

func (s ZapSinkRabbitMQ) Close() error {
	return s.closePublisher()
}

func (s ZapSinkRabbitMQ) closePublisher() error {
	if s.con == nil || s.con.IsClosed() {
		return nil
	}
	return s.con.Close()
}

func (s ZapSinkRabbitMQ) Write(p []byte) (n int, err error) {
	//不科学
	var m map[string]interface{}
	if err = json.Unmarshal(p, &m); err != nil {
		return
	}
	key, has := m[s.levelKey]
	if !has {
		err = errors.New("log lost level key")
		return
	}
	if key == nil || key.(string) == "" {
		err = errors.New("invalid log level key")
		return
	}
	err = s.ch.Publish(
		s.exchangeName, // exchange
		key.(string),   // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentEncoding: "utf8",
			ContentType:     "application/json",
			Body:            p,
		})
	return len(p), nil
}

func NewZapSinkRabbitMQ(url *url.URL) (Sink zap.Sink, err error) {
	s := &ZapSinkRabbitMQ{}
	s.exchangeName = url.Query().Get("exchange")
	if s.exchangeName == "" {
		err = errors.New("lost exchange name")
		return
	}
	s.levelKey = url.Query().Get("key")
	if s.levelKey == "" {
		err = errors.New("lost log level key name")
		return
	}

	s.errorChan = make(chan *amqp.Error, 1)
	if s.con, err = amqp.Dial(url.String()); err != nil {
		return
	}
	s.con.NotifyClose(s.errorChan)
	if s.ch, err = s.con.Channel(); err != nil {
		s.closePublisher()
		return
	}
	err = s.ch.ExchangeDeclare(
		s.exchangeName,      // name
		amqp.ExchangeDirect, // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		s.closePublisher()
		return
	}
	return s, nil
}

type ZapSinkMongo struct {
}

func (s ZapSinkMongo) Sync() error { return nil }

func (s ZapSinkMongo) Close() error { return nil }

func (s ZapSinkMongo) Write(p []byte) (n int, err error) {
	println("mongo->", string(p))
	return len(p), nil
}

func NewZapSinkMongo(url *url.URL) (sink zap.Sink, err error) {
	println(url.String())
	s := &ZapSinkMongo{}
	return s, nil
}

type ZapSinkXorm struct {
}

func (s ZapSinkXorm) Sync() error { return nil }

func (s ZapSinkXorm) Close() error { return nil }

func (s ZapSinkXorm) Write(p []byte) (n int, err error) { return }

func NewZapSinkXorm(url *url.URL) (sink zap.Sink, err error) { return }
