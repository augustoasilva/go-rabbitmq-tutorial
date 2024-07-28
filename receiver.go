package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, connErr := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if connErr != nil {
		// Este erro representa que falhamos ao abrir uma conexão com o RabbitMQ
		log.Panicf("[msg:%s][error:%s]", "failed to connect to RabbitMQ", connErr)
	}
	defer conn.Close()

	ch, chErr := conn.Channel()
	if chErr != nil {
		// Este erro representa que falhamos ao abrir um canal de comunicação com RabbitMQ
		log.Panicf("[msg:%s][error:%s]", "failed to open a channel", chErr)
	}
	defer ch.Close()

	q, qErr := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if qErr != nil {
		// Este erro representa que falhamos ao declarar uma fila nova no RabbitMQ
		log.Panicf("[msg:%s][error:%s]", "failed to declare a queue", qErr)
	}

	consumer, consumerErr := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if consumerErr != nil {
		// Este erro representa que falhamos ao declarar um consumidor do RabbitMQ
		log.Panicf("[msg:%s][error:%s]", "failed to declare a consumer", consumerErr)
	}

	var forever chan struct{}

	go func() {
		for message := range consumer {
			log.Printf("[msg:received a message][message.body:%s]", message.Body)
		}
	}()

	log.Printf("Waiting for messages. To exit press CTRL+C")
	<-forever
}
