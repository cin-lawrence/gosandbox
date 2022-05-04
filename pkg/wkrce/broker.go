package worker

import (
        "encoding/json"
        "time"

	"github.com/streadway/amqp"
        gc "github.com/gocelery/gocelery"
)

type AMQPCeleryBroker struct {
        *gc.AMQPCeleryBroker
}

func (b *AMQPCeleryBroker) SendCeleryMessage(message *gc.CeleryMessage) error {
	taskMessage := message.GetTaskMessage()
	queueName := "celery"
	_, err := b.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // autoDelete
		false,     // exclusive
		false,     // noWait
		nil,       // args
	)
	if err != nil {
		return err
	}
	err = b.ExchangeDeclare(
		"default",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	resBytes, err := json.Marshal(taskMessage)
	if err != nil {
		return err
	}

	publishMessage := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "application/json",
		Body:         resBytes,
	}

	return b.Publish(
		"",
		queueName,
		false,
		false,
		publishMessage,
	)
}

func NewAMQPCeleryBroker(host string) *AMQPCeleryBroker {
        conn, channel := gc.NewAMQPConnection(host)
	broker := &AMQPCeleryBroker{
                &gc.AMQPCeleryBroker{
                        Channel:    channel,
                        Connection: conn,
                        Exchange:   &gc.AMQPExchange{
                                Name: "default",
                                Type: "direct",
                                Durable: true,
                                AutoDelete: false,
                        },
                        Queue:      &gc.AMQPQueue{
                                Name: "celery",
                                Durable: true,
                                AutoDelete: false,
                        },
                        Rate:   1,
                },
	}
	if err := broker.CreateExchange(); err != nil {
		panic(err)
	}
	if err := broker.CreateQueue(); err != nil {
		panic(err)
	}
	if err := broker.Qos(broker.Rate, 0, false); err != nil {
		panic(err)
	}
	if err := broker.StartConsumingChannel(); err != nil {
		panic(err)
	}
	return broker
}
