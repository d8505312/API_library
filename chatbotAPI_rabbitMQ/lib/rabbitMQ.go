package lib

import (
	"fmt"

	"github.com/streadway/amqp"
)

var MQ *amqp.Connection
var MQ2 *amqp.Connection

func InitMQ() {
	MQ = ConnMQ(
		Config.GetString("rabbitMQ.user"),
		Config.GetString("rabbitMQ.pass"),
		Config.GetString("rabbitMQ.ip"),
		Config.GetString("rabbitMQ.port"),
		Config.GetString("rabbitMQ.name"),
	)
}

func InitMQ2() {
	MQ2 = ConnMQ(
		Config.GetString("rabbitMQ.user"),
		Config.GetString("rabbitMQ.pass"),
		Config.GetString("rabbitMQ.ip"),
		Config.GetString("rabbitMQ.port"),
		"comcloud",
	)
}

//連接 RabbitMQ
func ConnMQ(account, password, host, port, name string) *amqp.Connection {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%v:%v@%v:%v/%v", account, password, host, port, name))
	FailOnError(err, "Failed to connect to RabbitMQ")
	return conn
}

func GetMQ() *amqp.Connection {
	return MQ
}

func GetMQ2() *amqp.Connection {
	return MQ2
}

func FailOnError(err error, msg string) {
	if err != nil {
		Log.Fatalf("%s: %s", msg, err)
	}
}
