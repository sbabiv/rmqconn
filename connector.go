package rmqconn

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/streadway/amqp"
)


//Conn connect instance
type Conn struct {
	sync.Mutex
	connection Conner

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
	// ErrConnectionUnavailable remote server is unavailable or network is down
	ErrConnectionUnavailable = errors.New("connection unavailable")
	// ErrClsoe connection was closed earlier
	ErrClose = errors.New("connection was closed earlier")
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
	instance := &Conn{url: url, dial: dial, done: make(chan struct{})}
	err := instance.conn()
	go func() {
		instance.recon()
	}()

	return instance, err
}

// IsConnected return connection status
func (conn *Conn) IsConnected() bool {
	return atomic.LoadInt32(&conn.isConnected) == connected
}

// Close re-conn attempts
func (conn *Conn) Close() error {
	conn.Lock()
	defer conn.Unlock()

	if atomic.LoadInt32(&conn.isConnected) == closed {
		return ErrClose
	}

	atomic.StoreInt32(&conn.isConnected, closed)
	close(conn.done)

	if conn.connection != nil {
		return conn.connection.Close()
	}

	return nil
}

func (conn *Conn) GetChannel() (*amqp.Channel, error) {
	if conn.IsConnected() {
		return conn.connection.GetChannel()
	}
	return nil, ErrConnectionUnavailable
}

func (conn *Conn) recon() {
	for {
		for {
			if atomic.LoadInt32(&conn.isConnected) == connected {
				break
			}

			if atomic.LoadInt32(&conn.isConnected) == closed {
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
			if atomic.LoadInt32(&conn.isConnected) != closed {
				atomic.StoreInt32(&conn.isConnected, disconnected)
			}
		}
	}
}

func (conn *Conn) conn() error {
	conn.Lock()
	defer conn.Unlock()

	c, err := conn.dial(conn.url)
	if err != nil {
		return err
	}

	conn.connection = c
	conn.notifyClose = make(chan *amqp.Error)
	conn.connection.NotifyClose(conn.notifyClose)

	atomic.StoreInt32(&conn.isConnected, connected)

	return nil
}