[![Build Status](https://travis-ci.org/sbabiv/rmqconn.svg?branch=master)](https://travis-ci.org/sbabiv/rmqconn)
[![cover.run](https://cover.run/go/github.com/sbabiv/rmqconn.svg?style=flat&tag=golang-1.10)](https://cover.run/go?tag=golang-1.10&repo=github.com%2Fsbabiv%2Frmqconn)
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

Or, using dep:

``` sh
dep ensure -add github.com/sbabiv/rmqconn
```


#### 2. use it

```Go
conn, err := rmqconn.Open("amqp://", rmqconn.Dial)
defer conn.Close()

if err != nil {
  return
}

err = conn.Do(func(ch *amqp.Channel) error {
  return ch.Publish("", "queueName", false, false, amqp.Publishing{
    Body: []byte("hello wolrd"),
  })
})
  ```
## Licence
[MIT](https://opensource.org/licenses/MIT)

## Author 
Babiv Sergey
