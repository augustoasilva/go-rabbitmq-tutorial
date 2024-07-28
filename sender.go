// sender.go
package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Como não passamos nenhuma configuração no arquivo do docker-compose.yml
	// nossa url de conexão será a seguinte abaixo.
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
		"hello", // nome da fila
		false,   // atributo para dizer se a fila é duravel
		false,   // se a fila é auto deletavel quando não usada
		false,   // atributo para dizer se a fila é exculsiva ou não
		false,   // atributo para dizer se a fila deve esperar ou não
		nil,     // argumentos de configuração da fila
	)
	if qErr != nil {
		// Este erro representa que falhamos ao declarar uma fila nova no RabbitMQ
		log.Panicf("[msg:%s][error:%s]", "failed to declare a queue", qErr)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Olá RabbitMQ! Minha primeira mensagem a ser enviada"
	publishErr := ch.PublishWithContext(ctx,
		"",     // "exchange" a ser usada no RabbitMQ
		q.Name, // routing key (chave de roteamento)  a ser usada no RabbitMQ
		false,  // atributo que indica se é mandatório o envio
		false,  // atributo que indicai se é imediato o envio
		amqp.Publishing{
			ContentType: "text/plain", // tipo da mensagem a ser enviada
			Body:        []byte(body),
		})
	if publishErr != nil {
		// Este erro representa que falhamos ao enviar uma mensagem para o RabbitMQ
		log.Panicf("[msg:%s][error:%s]", "failed to publish a message", qErr)
	}

	log.Printf("[msg:message sent to rabbitmq][body:%s]\n", body)
}
