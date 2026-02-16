package rabbitmq

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/quickybrains/bankapp/internal/domain"
	"github.com/quickybrains/bankapp/internal/ports/log"
	"github.com/quickybrains/bankapp/internal/ports/notifier"
	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitMQConsumer struct {
	amqpCh <-chan amqp.Delivery

	wg     sync.WaitGroup
	stopCh chan struct{}

	config *notifier.Config
	log    log.Logger
}

func StartNewConsumer(conf *notifier.Config, log log.Logger) (*rabbitMQConsumer, error) {
	consumer := &rabbitMQConsumer{
		config: conf,
		log:    log,

		stopCh: make(chan struct{}, 1),
	}

	err := consumer.init()
	if err != nil {
		return nil, fmt.Errorf("notifier init, err=%w", err)
	}

	consumer.wg.Add(1)
	go consumer.consumeLoop()

	return consumer, nil
}

func (consumer *rabbitMQConsumer) Stop() error {
	consumer.stopCh <- struct{}{}

	consumer.wg.Wait()

	return nil
}

func (consumer *rabbitMQConsumer) init() error {
	addr := consumer.config.DSN()

	consumer.log.With("addr", addr).Info().Print("Connecting to rabbitmq")

	conn, err := amqp.Dial(addr)
	if err != nil {
		return fmt.Errorf("amqp dial, addr=%s, err=%w", addr, err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("amqp conn channel, err=%w", err)
	}

	err = ch.ExchangeDeclare(
		exchangeName,
		amqp.ExchangeFanout,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("amqp channel exchange declare, err=%w", err)
	}

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("amqp channel queue declare, err=%w", err)
	}

	err = ch.QueueBind(
		q.Name,       // queue name (generated)
		"",           // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("amqp channel queue bind, err=%w", err)
	}

	consumer.amqpCh, err = ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	consumer.log.With("addr", addr).Info().Print("Consumer started")

	return nil
}

func (consumer *rabbitMQConsumer) consumeLoop() {
	defer consumer.wg.Done()

	for {
		select {
		case delivery := <-consumer.amqpCh:
			var msg domain.AccountBalanceChangedMsg
			err := json.Unmarshal(delivery.Body, &msg)
			if err != nil {
				consumer.log.Err(err).Print("json unmarshal")

				continue
			}
			consumer.log.Info().Print(msg.String())
		case <-consumer.stopCh:
			consumer.log.Info().Print("Stopping consumer on signal")
			return
		}
	}
}
