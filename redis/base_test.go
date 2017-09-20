package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"math/rand"
	"flag"
)

var tstRedisHost string
var tstRedisPass string
var tstRedisPort int
var tstRedisDb int

func init() {
	flag.StringVar(&tstRedisHost, "redis.host", "127.0.0.1", "the redis hostname/ip")
	flag.IntVar(&tstRedisPort, "redis.port", 6379, "redis port")
	flag.IntVar(&tstRedisDb, "redis.db", 5, "redis db")
	flag.StringVar(&tstRedisPass, "redis.pass", "", "redis password")
	flag.Parse()
}

func makeProducer(opt *redis.Options, queueName string) (func(msg []byte), func()) {
	client := redis.NewClient(opt)
	producerFunc := func(msg []byte) {
		client.RPush(queueName, msg)
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

func makeTestOpts() (*redis.Options, string) {
	opt := &redis.Options{
		Addr:       fmt.Sprintf("%s:%d", tstRedisHost, tstRedisPort),
		Password:   tstRedisPass,
		DB:         tstRedisDb,
		MaxRetries: 2,
		PoolSize:   20,
	}
	queue := randStringRunes(10)
	return opt, queue
}
