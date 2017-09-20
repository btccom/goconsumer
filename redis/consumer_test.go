package redis

import (
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	opts, queueName := makeTestOpts()

	consumer := New(opts, queueName)
	assert.NotNil(t, consumer)
	assert.IsType(t, Consumer{}, *consumer)
	assert.NotNil(t, consumer.redis)
	assert.IsType(t, redis.Client{}, *consumer.redis)
	assert.Equal(t, queueName, consumer.redisKey)
}

func TestConsumeInternallyTimesOut(t *testing.T) {
	opts, queueName := makeTestOpts()
	consumer := New(opts, queueName)

	waitSeconds := 3

	start := time.Now()
	doneCount := make(chan int)

	go func(doneCount chan int, waitSeconds int) {
		for i := 0; i < waitSeconds; i++ {
			val, err := consumer.consume()
			assert.NoError(t, err)
			assert.NotNil(t, val)
			assert.True(t, len(val) == 0)
		}
		doneCount <- waitSeconds
	}(doneCount, waitSeconds)

	done := <-doneCount

	assert.Equal(t, waitSeconds, done)
	assert.True(t, time.Since(start) > time.Second*time.Duration(waitSeconds))
}

func TestRedisConsumer(t *testing.T) {
	input := make([][]byte, 3)
	input[0] = []byte{0x41}
	input[1] = []byte{0x41, 0x42}
	input[2] = []byte{0x41, 0x42, 0x43}

	opts, queueName := makeTestOpts()
	producer, closer := makeProducer(opts, queueName)
	defer closer()

	consumer := New(opts, queueName)
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
		if len(received) == len(input) {
			consumer.Close()
		}
	}
}

func TestIsUnBuffered(t *testing.T) {

	opts, queueName := makeTestOpts()
	producer, closer := makeProducer(opts, queueName)
	defer closer()

	consumer := New(opts, queueName)

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
			consumer.consume()
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
