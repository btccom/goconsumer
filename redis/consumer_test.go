package redis

import (
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	opts, queueName := makeTestOpts("127.0.0.1", 6979, "")

	consumer := New(opts, queueName)
	assert.NotNil(t, consumer)
	assert.IsType(t, Consumer{}, *consumer)
	assert.NotNil(t, consumer.redis)
	assert.IsType(t, redis.Client{}, *consumer.redis)
	assert.Equal(t, queueName, consumer.redisKey)
}

func TestMockConsumer(t *testing.T) {
	input := make([][]byte, 3)
	input[0] = []byte{0x41}
	input[1] = []byte{0x41, 0x42}
	input[2] = []byte{0x41, 0x42, 0x43}

	opt, queueName := makeTestOpts("127.0.0.1", 6979, "")
	producer, closer := makeProducer(opt, queueName)
	defer closer()

	consumer := New(opt, queueName)
	go func(producer func([]byte), list ...[]byte) {
		l := len(list)
		for i := 0; i < l; i++ {
			producer(list[i])
		}
	}(producer, input...)

	ch := consumer.Channel()
	received := make([][]byte, 0)
	for msg := range ch {
		received = append(received, msg)
		checkAgainst := input[len(received)-1]
		assert.Equal(t, checkAgainst, msg)
	}

	// Channel is closed
	v := <-ch
	assert.Nil(t, v)
}

func TestIsUnBuffered(t *testing.T) {

	opt, queueName := makeTestOpts("127.0.0.1", 6979, "")
	producer, closer := makeProducer(opt, queueName)
	defer closer()

	consumer := New(opt, queueName)

	start := time.Now()
	done := make(chan bool)

	numValues := 3
	delay := time.Millisecond * 250

	go func(producer func([]byte), numValues int) {
		for i := 0; i < numValues; i++ {
			producer([]byte{0x41})
		}
	}(producer, numValues)

	go func(consumer *Consumer, doneChan chan bool, delay time.Duration) {
		for i := 0; i < numValues; i++ {
			time.Sleep(delay)
			consumer.Consume()
		}
		done <- true
	}(consumer, done, delay)

	<-done
	diff := time.Since(start)

	expectDelay := time.Duration(0)
	for i := 0; i < numValues; i++ {
		expectDelay += delay
	}

	assert.True(t, diff > expectDelay)
}
