package rmqconn

import (
	"errors"
	"github.com/streadway/amqp"
	"sync"
	"sync/atomic"
	"time"
)

// Connecter interface for connection instance
type Connecter interface {
	Do(f func(ch *amqp.Channel) error) error
	IsСonnected() bool
	Close() error
}

//Conn connect instance
type Conn struct {
	sync.Mutex
	connection Conner
	channel    Channer

	notifyClose chan *amqp.Error

	isConnected int32
	attempts    int8
	url         string

	dial func(string) (Conner, error)
}

const (
	disconnected int32 = 0
	connected    int32 = 1
	closed       int32 = -1

	maxAttempts int8 = 5
)

var (
	//ErrConnectionUnavailable remote server is unavailable or network is down
	ErrConnectionUnavailable = errors.New("connection unavailable")
)

// Dial function wrapper amqp.Dial
// For use amqp.DialConfig or amqp.DialTLS implement your func
func Dial(url string) (Conner, error) {
	c, err := amqp.Dial(url)

	return &connWrapper{conn: c}, err
}

// Creates connection and trying to connect
// use Close method to stop connection attempts
func Open(url string, dial func(string) (Conner, error)) (Connecter, error) {
	instance := &Conn{url: url, dial: dial}
	err := instance.conn()
	go func() {
		instance.recon()
	}()

	return instance, err
}

// Do return *amqp.Channel if instance connected
func (conn *Conn) Do(f func(ch *amqp.Channel) error) error {
	if atomic.LoadInt32(&conn.isConnected) != connected || conn.channel == nil {
		return ErrConnectionUnavailable
	}

	return f(conn.channel.GetChannel())
}

// IsСonnected return connection status
func (conn *Conn) IsСonnected() bool {
	return atomic.LoadInt32(&conn.isConnected) == connected
}

// Close re-conn attempts
func (conn *Conn) Close() error {
	var err error

	isAlive := atomic.LoadInt32(&conn.isConnected)
	if isAlive == closed {
		return err
	}

	atomic.StoreInt32(&conn.isConnected, closed)

	if conn.channel != nil {
		err = conn.channel.Close()
		if err != nil {
			return err
		}
	}

	if conn.connection != nil {
		err = conn.connection.Close()
	}

	return err
}

func (conn *Conn) recon() {
	for {
		for {
			isAlive := atomic.LoadInt32(&conn.isConnected)
			if isAlive == connected {
				break
			}

			if isAlive == closed {
				return
			}

			if err := conn.conn(); err != nil {
				if conn.attempts < maxAttempts {
					conn.attempts++
				}
			} else {
				conn.attempts = 0
			}

			time.Sleep(time.Second * time.Duration(conn.attempts))
		}

		select {
		case <-conn.notifyClose:
			if atomic.LoadInt32(&conn.isConnected) == closed {
				return
			}
			atomic.StoreInt32(&conn.isConnected, disconnected)
		}
	}
}

func (conn *Conn) conn() error {
	c, err := conn.dial(conn.url)
	if err != nil {
		return err
	}

	ch, err := c.Channel()
	if err != nil {
		return err
	}

	conn.connection = c
	conn.channel = ch
	conn.notifyClose = make(chan *amqp.Error)
	conn.channel.NotifyClose(conn.notifyClose)

	atomic.StoreInt32(&conn.isConnected, connected)

	return nil
}
