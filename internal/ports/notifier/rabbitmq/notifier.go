package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/quickybrains/bankapp/internal/domain"
	"github.com/quickybrains/bankapp/internal/ports/log"
	"github.com/quickybrains/bankapp/internal/ports/notifier"

	amqp "github.com/rabbitmq/amqp091-go"
)

const exchangeName = "balance_changed"

type rabbitMQNotifier struct {
	amqpCh *amqp.Channel

	config *notifier.Config
	log    log.Logger
}

func NewNotifier(conf *notifier.Config, log log.Logger) (*rabbitMQNotifier, error) {
	notifier := &rabbitMQNotifier{
		config: conf,
		log:    log,
	}

	err := notifier.init()
	if err != nil {
		return nil, fmt.Errorf("notifier init, err=%w", err)
	}

	return notifier, nil
}

func (notifier *rabbitMQNotifier) NotifyBalanceChanged(ctx context.Context, message *domain.AccountBalanceChangedMsg) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("json marshal, err=%w", err)
	}

	err = notifier.amqpCh.PublishWithContext(ctx,
		exchangeName, // exchange
		"",           // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json:",
			Body:        []byte(data),
		})

	return nil
}

func (notifier *rabbitMQNotifier) init() error {
	amqpURL := notifier.config.DSN()

	notifier.log.With("amqp_url", amqpURL).Info().Print("Connecting to rabbitmq")

	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return fmt.Errorf("amqp dial, url=%s, err=%w", amqpURL, err)
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

	notifier.log.
		With("exchange", exchangeName).
		Info().
		Print("Connection established, exhange declared")

	return nil
}
