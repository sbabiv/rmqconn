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
	Open(url string, dial func(string) (Conner, error)) (*Conn, error)
	Do(f func(ch *amqp.Channel) error) error
	IsСonnected() bool
	Close() error
}

//Conn connect instance
type Conn struct {
	connection Conner
	channel    Channer

	done        chan struct{}
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
	instance *Conn
	once     sync.Once
	openErr  error

	//ErrNoConnection remote server is unavailable or network is down
	ErrNoConnection = errors.New("no connection")
)

// Dial function wrapper amqp.Dial
// For use amqp.DialConfig or amqp.DialTLS implement your func
func Dial(url string) (Conner, error) {
	c, err := amqp.Dial(url)
	return &connWrapper{conn: c}, err
}

// Open safe for concurrent use by multiple goroutines. Creates a single connection once (singleton). And trying to connect
// use Close method to stop connection attempts
func Open(url string, dial func(string) (Conner, error)) (*Conn, error) {
	once.Do(func() {
		instance = &Conn{url: url, dial: dial}
		openErr = instance.conn()
		go func() {
			instance.recon()
		}()
	})

	return instance, openErr
}

// Do return *amqp.Channel if instance connected
func (conn *Conn) Do(f func(ch *amqp.Channel) error) error {
	if atomic.LoadInt32(&conn.isConnected) != connected || conn.channel == nil {
		return ErrNoConnection
	}

	return f(conn.channel.GetChannel())
}

// IsСonnected return connection status
func (conn *Conn) IsСonnected() bool {
	return atomic.LoadInt32(&conn.isConnected) == connected
}

// Close re-conn attempts
func (conn *Conn) Close() error {
	isAlive := atomic.LoadInt32(&conn.isConnected)
	if isAlive == closed || isAlive == disconnected {
		return nil
	}

	atomic.StoreInt32(&conn.isConnected, closed)
	close(conn.done)

	err := conn.channel.Close()
	if err != nil {
		return err
	}
	err = conn.connection.Close()
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
		case <-conn.done:
			return
		case <-conn.notifyClose:
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
	conn.done = make(chan struct{})
	conn.channel.NotifyClose(conn.notifyClose)

	atomic.StoreInt32(&conn.isConnected, connected)

	return nil
}
