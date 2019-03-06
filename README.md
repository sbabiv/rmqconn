# rmqconn
RabbitMQ Reconnection for Golang

[![Build Status](https://travis-ci.org/sbabiv/rmqconn.svg?branch=master)](https://travis-ci.org/sbabiv/rmqconn)
[![Coverage Status](https://coveralls.io/repos/github/sbabiv/rmqconn/badge.svg?branch=master)](https://coveralls.io/github/sbabiv/rmqconn?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/sbabiv/rmqconn)](https://goreportcard.com/report/github.com/sbabiv/rmqconn)
[![GoDoc](https://godoc.org/github.com/sbabiv/rmqconn?status.svg)](https://godoc.org/github.com/sbabiv/rmqconn)

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
