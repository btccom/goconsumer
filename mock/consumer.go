package mock

func New() *Consumer {
	return &Consumer{
		srcChan: make(chan []byte),
		msgs:    make(chan []byte),
		quit:    make(chan bool),
		open:    false,
	}
}

type Consumer struct {
	srcChan chan []byte
	msgs    chan []byte
	quit    chan bool
	open    bool
}

func (c *Consumer) Produce(msg []byte) {
	c.srcChan <- msg
}

func (c *Consumer) Consume() ([]byte, error) {
	msg := <-c.srcChan
	return msg, nil
}

func (c *Consumer) Channel() chan []byte {
	done := make(chan bool)
	go func() {
		done <- true
		for {
			select {
			case msg := <-c.srcChan:
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
