package goconsumer

type Consumer interface {
	Consume() ([]byte, error)
	Channel() chan []byte
	Close()
}
