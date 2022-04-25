package controller

import (
	"encoding/json"
	"fmt"

	. "chatbotAPI/lib"

	"github.com/streadway/amqp"
)

//發送訊息至MQ
func SendToMQ(replytrsend []byte) {

	data := replytrsend
	fmt.Printf("data:\n%s\n", json.RawMessage(data))

	mq := GetMQ2()

	ch, err := mq.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"message", // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		})
	FailOnError(err, "Failed to publish a message")
}
