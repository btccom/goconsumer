package goconsumer

type Consumer interface {
	Channel() chan []byte
	Close()
}
