package redis

import (
	"testing"
	"fmt"
)

func BenchmarkConsumer(b *testing.B) {
	opts, queueName := makeTestOpts("127.0.0.1", 6979, "")
	consumer := New(opts, queueName)
	producer, closer := makeProducer(opts, queueName)
	defer closer()

	go func(producer func ([]byte)) {
		toSend := []byte{0x42}
		for n := 0; n < b.N; n++ {
			producer(toSend)
		}
	}(producer)

	ch := consumer.Channel()
	doneChan := make(chan int)

	go func(doneChan chan int) {
		result := 0
		for n := 0; n < b.N; n++ {
			<- ch
			result += 1
		}

		doneChan <- result
	}(doneChan)

	result := <- doneChan
	fmt.Printf("Benchmark covered %d receives\n", result)
}