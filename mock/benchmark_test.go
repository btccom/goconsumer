package mock

import (
	"fmt"
	"testing"
)

func BenchmarkConsumer(b *testing.B) {
	consumer := New()
	go func() {
		toSend := []byte{0x42}
		for n := 0; n < b.N; n++ {
			consumer.Produce(toSend)
		}
		consumer.Close()
	}()

	ch := consumer.Channel()
	doneChan := make(chan int)

	go func(doneChan chan int) {
		result := 0
		for n := 0; n < b.N; n++ {
			<-ch
			result += 1
		}

		doneChan <- result
	}(doneChan)

	result := <-doneChan
	fmt.Printf("Benchmark covered %d receives\n", result)
}
