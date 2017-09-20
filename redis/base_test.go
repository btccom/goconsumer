package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"math/rand"
	"time"
)

func makeProducer(opt *redis.Options, queueName string) (func(msg []byte), func()) {
	client := redis.NewClient(opt)
	producerFunc := func(msg []byte) {
		client.RPush(queueName, msg)
		time.Sleep(time.Millisecond * 200)
	}
	closeFunc := func() {
		client.Close()
	}

	return producerFunc, closeFunc
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func makeTestOpts(host string, port int, pass string) (*redis.Options, string) {
	opt := &redis.Options{
		Addr:       fmt.Sprintf("%s:%d", host, port),
		Password:   pass,
		DB:         5,
		MaxRetries: 2,
		PoolSize:   20,
	}
	queue := randStringRunes(10)
	return opt, queue
}
