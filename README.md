[![Build Status](https://travis-ci.org/sbabiv/rmqconn.svg?branch=master)](https://travis-ci.org/sbabiv/rmqconn)
[![Coverage Status](https://coveralls.io/repos/github/sbabiv/rmqconn/badge.svg?branch=master)](https://coveralls.io/github/sbabiv/rmqconn?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/sbabiv/rmqconn)](https://goreportcard.com/report/github.com/sbabiv/rmqconn)
[![GoDoc](https://godoc.org/github.com/sbabiv/rmqconn?status.svg)](https://godoc.org/github.com/sbabiv/rmqconn)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#messaging)

# rmqconn
RabbitMQ Reconnection for Golang    

Wrapper over `amqp.Connection` and `amqp.Dial`. Allowing to do a reconnection when the connection is broken before forcing the call to the Close () method to be closed

Use the default method `func Dial (url string) (Conner, error)` to connect to the server.
You can implement your connection function and pass it to `rmqconn.Open("", customFunc)`

## Getting started

#### 1. install

``` sh
go get -u github.com/sbabiv/rmqconn
```

#### 2. use it

```Go
conn, err := rmqconn.Open("amqp://usr:pwd@host:5672", rmqconn.Dial)
defer conn.Close()

if err != nil {
    return
}

if conn.IsConnected() {
    ch, err := conn.GetChannel()
    if err != nil {
        return
    }
    defer ch.Close()

    err = ch.Publish("", "queueName", false, false, amqp.Publishing{
        Body: []byte("hello wolrd"),
    })
}
  ```
## Licence
[MIT](https://opensource.org/licenses/MIT)

## Author 
Babiv Sergey
