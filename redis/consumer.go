package redis

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"strings"
	"time"
)

func New(cfg *redis.Options, channel string) *Consumer {
	return &Consumer{
		redis:    redis.NewClient(cfg),
		redisKey: channel,
		msgs:     make(chan []byte),
		quit:     make(chan bool),
	}
}

type Consumer struct {
	redis    *redis.Client
	redisKey string
	msgs     chan []byte
	quit     chan bool
}

func (c *Consumer) consume() ([]byte, error) {
	timeout := 1 * time.Second

	msg := c.redis.BRPop(timeout, c.redisKey)
	res, err := msg.Result()

	var ret string
	if err != nil {
		ret, err = c.handleError(err)
		if err != nil {
			return nil, err
		}
	} else {
		ret = res[1]
	}

	return []byte(ret), nil
}

func (c *Consumer) handleError(err error) (string, error) {
	// `redis: nil` indicates we just reached the set blocking timeout
	if err.Error() == "redis: nil" {
		return "", nil

		// i/o timeout err indicates we just reached the set blocking timeout
	} else if strings.HasSuffix(err.Error(), "i/o timeout") {
		return "", nil

	} else {
		return "", errors.Wrapf(err, "brpop failed")
	}
}

func (c *Consumer) Channel() chan []byte {
	done := make(chan bool)
	go func() {
		done <- true
		for {
			select {
			case <-c.quit:
				close(c.msgs)
				return
			default:
				msg, err := c.consume()
				if err != nil {
					close(c.msgs)
					return
				}
				if len(msg) > 0 {
					c.msgs <- msg
				}
			}
		}
	}()
	<-done
	return c.msgs
}

func (c *Consumer) Close() {
	c.quit <- true
}
