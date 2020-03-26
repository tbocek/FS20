package main

import (
	"log"

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

	msgs, err := ch.Consume(
		q.Name, "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
	}
}
