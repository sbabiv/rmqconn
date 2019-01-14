package rmqconn

import (
	"fmt"
	"github.com/streadway/amqp"
)

// Channer interface for wrapper Channel
type Channer interface {
	Close() error
	GetChannel() *amqp.Channel
	NotifyClose(c chan *amqp.Error) chan *amqp.Error
}

type chann struct {
	c *amqp.Channel
}

func (ch *chann) Close() error {
	if ch.c == nil {
		return fmt.Errorf("channel nil")
	}
	return ch.c.Close()
}

func (ch *chann) GetChannel() *amqp.Channel {
	return ch.c
}

func (ch *chann) NotifyClose(c chan *amqp.Error) chan *amqp.Error {
	return ch.c.NotifyClose(c)
}

// Conner interface for wrapper Connection
type Conner interface {
	Close() error
	Channel() (Channer, error)
}

type connWrapper struct {
	conn *amqp.Connection
}

func (c *connWrapper) Close() error {
	return c.conn.Close()
}

func (c *connWrapper) Channel() (Channer, error) {
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, err
	}

	return &chann{c: ch}, nil
}
