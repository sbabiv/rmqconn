package rmqconn

import (
	"fmt"
	"github.com/streadway/amqp"
	"sync/atomic"
	"testing"
	"time"
)

var (
	attempts int32
)

type connMock struct {
	с chan *amqp.Error
}

func (c *connMock) Close() error {
	close(c.с)
	return nil
}

func (c *connMock) NotifyClose(ch chan *amqp.Error) chan *amqp.Error {
	c.с = ch
	return nil
}

func (c *connMock) GetChannel() (*amqp.Channel, error) {
	return new(amqp.Channel), nil
}

func (c *connMock) Disconnect() {
	close(c.с)
}

func DialMock(url string) (Conner, error) {
	defer atomic.AddInt32(&attempts, 1)

	if atomic.LoadInt32(&attempts) == 0 {
		conn := &connMock{make(chan *amqp.Error)}
		go func() {
			time.Sleep(time.Millisecond * 100)
			conn.Disconnect()
		}()
		return conn, nil
	}

	if atomic.LoadInt32(&attempts) <= 1 {
		return nil, fmt.Errorf("dial err")
	} else {
		return &connMock{make(chan *amqp.Error)}, nil
	}
}

func TestReconn(t *testing.T) {
	c, err := Open("amqp://host", DialMock)
	if err != nil {
		t.Fail()
	}
	if !c.IsConnected() {
		t.Fail()
	} else {
		_, err := c.GetChannel()
		if err != nil {
			t.Fail()
		}
	}

	time.Sleep(time.Second * 3)

	c.Close()
	c.Close()

	_, err = c.GetChannel()
	if err == nil {
		t.Fail()
	}
}

func TestDial(t *testing.T) {
	c, err := Open("amqp://host", Dial)
	if err == nil {
		t.Fail()
	}
	c.Close()
}