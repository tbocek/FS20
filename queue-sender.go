package main

import (
	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("storage", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	body := "Hallo HSR 12345"
	err = ch.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})
	if err != nil {
		panic(err)
	}
}
