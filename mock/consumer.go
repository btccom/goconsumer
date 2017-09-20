package mock

func New() *Consumer {
	return &Consumer{
		SrcChan: make(chan []byte),
		msgs:    make(chan []byte),
		quit:    make(chan bool),
		open:    false,
	}
}

type Consumer struct {
	SrcChan chan []byte
	msgs    chan []byte
	quit    chan bool
	open    bool
}

func (c *Consumer) Produce(msg []byte) {
	c.SrcChan <- msg
}

func (c *Consumer) Consume() ([]byte, error) {
	msg := <-c.SrcChan
	return msg, nil
}

func (c *Consumer) Channel() chan []byte {
	done := make(chan bool)
	go func() {
		done <- true
		for {
			select {
			case msg := <-c.SrcChan:
				c.msgs <- msg
			case <- c.quit:
				close(c.msgs)
			}
		}
	}()
	<-done
	return c.msgs
}

func (c *Consumer) Close() {
	c.quit <- true
}
