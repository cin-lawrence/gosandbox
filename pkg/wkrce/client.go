package worker

import (
	"time"

	gc "github.com/gocelery/gocelery"
	log "github.com/sirupsen/logrus"

	"github.com/cin-lawrence/gosandbox/pkg/config"
)

func NewWorkerClient() *gc.CeleryClient {
	numAttempts := 0
	for {
		client, err := gc.NewCeleryClient(
			NewAMQPCeleryBroker(config.Config.BrokerURI),
			gc.NewAMQPCeleryBackend(config.Config.BrokerURI),
			1,
		)
		if err != nil {
			if numAttempts < 30 {
				log.Warn("Can't connect to AMQP, retrying...")
				numAttempts += 1
				time.Sleep(2)
				continue
			}
			panic(err)
		}
		return client
	}
}
