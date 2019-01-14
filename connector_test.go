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

type channMock struct {
	с chan *amqp.Error
}

func (ch *channMock) Close() error {
	close(ch.с)
	return nil
}

func (ch *channMock) NotifyClose(c chan *amqp.Error) chan *amqp.Error {
	ch.с = c

	return nil
}

func (ch *channMock) GetChannel() *amqp.Channel {
	return new(amqp.Channel)
}

type connMock struct {
}

func (c *connMock) Close() error {
	return nil
}

func (c *connMock) Channel() (Channer, error) {
	if atomic.LoadInt32(&attempts) == 3 {
		return new(channMock), nil
	}
	return new(channMock), fmt.Errorf("get chan err")
}

func DialMock(url string) (Conner, error) {
	defer atomic.AddInt32(&attempts, 1)
	if atomic.LoadInt32(&attempts) <= 2 {
		return new(connMock), nil
	}
	return new(connMock), fmt.Errorf("dial err")
}

func TestConn(t *testing.T) {
	c, err := Open("amqp://", DialMock)
	if err == nil {
		t.Fail()
	}
	time.Sleep(time.Second * 3)

	err = c.Do(func(ch *amqp.Channel) error {
		return nil
	})

	if err != nil {
		t.Fail()
	}
	if !c.IsСonnected() {
		t.Fail()
	}

	c.Close()
	c.Do(func(ch *amqp.Channel) error {
		return nil
	})
	c.Close()

	c, _ = Open("amqp://", DialMock)
	c.Close()
}

func TestDial(t *testing.T) {
	c, err := Dial("amqp://")
	if c != nil && err == nil {
		t.Fail()
	}
}

func TestTypes(t *testing.T) {
	ch := &chann{c :new(amqp.Channel)}
	ch.GetChannel()
	c := make(chan *amqp.Error)
	ch.NotifyClose(c)
}