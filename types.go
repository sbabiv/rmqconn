package rmqconn

import (
	"github.com/streadway/amqp"
)

// Conner interface for wrapper Connection
type Conner interface {
	Close() error
	NotifyClose(c chan *amqp.Error) chan *amqp.Error
	GetChannel() (*amqp.Channel, error)
}

// Connecter interface for connection instance
type Connecter interface {
	GetChannel() (*amqp.Channel, error)
	IsConnected() bool
	Close() error
}

type connWrapper struct {
	conn *amqp.Connection
}

func (cw *connWrapper) Close() error {
	return cw.conn.Close()
}

func (cw *connWrapper) NotifyClose(c chan *amqp.Error) chan *amqp.Error {
	return cw.conn.NotifyClose(c)
}

func (cw *connWrapper) GetChannel() (*amqp.Channel, error) {
	return cw.conn.Channel()
}