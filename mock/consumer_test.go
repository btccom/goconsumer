package mock

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	consumer := New()
	assert.NotNil(t, consumer)
	assert.IsType(t, Consumer{}, *consumer)
}

func TestMockConsumer(t *testing.T) {
	input := make([][]byte, 3)
	input[0] = []byte{0x41}
	input[1] = []byte{0x41, 0x42}
	input[2] = []byte{0x41, 0x42, 0x43}

	consumer := New()
	go func(list ...[]byte) {
		l := len(list)
		for i := 0; i < l; i++ {
			consumer.Produce(list[i])
		}
		consumer.Close()
	}(input...)

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
	consumer := New()

	start := time.Now()
	done := make(chan bool)

	numValues := 3
	delay := time.Millisecond * 250

	go func(consumer *Consumer, doneChan chan bool, numValues int) {
		for i := 0; i < numValues; i++ {
			consumer.Produce([]byte{0x41})
		}
		doneChan <- true
	}(consumer, done, numValues)

	go func(consumer *Consumer, delay time.Duration) {
		for i := 0; i < numValues; i++ {
			time.Sleep(delay)
			<-consumer.SrcChan
		}
	}(consumer, delay)

	<-done
	diff := time.Since(start)

	expectDelay := time.Duration(0)
	for i := 0; i < numValues; i++ {
		expectDelay += delay
	}

	assert.True(t, diff > expectDelay)
}
