package nats

import (
	"context"
	"fmt"
	cenats "github.com/cloudevents/sdk-go/protocol/nats/v2"
	ce "github.com/cloudevents/sdk-go/v2"
	"log/slog"
	"kubercode-sso/config"
)

type MessageBroker interface {
	Send(ctx context.Context, message ce.Event) error
	Close(ctx context.Context) error
}

type CEMessageBroker struct {
	log    *slog.Logger
	cfg    *config.Config
	sender *cenats.Sender
	client ce.Client
}

func NewCEMessageBroker(config *config.Config, log *slog.Logger) *CEMessageBroker {
	sender, err := cenats.NewSender(config.NatsURL, config.Subject, cenats.NatsOptions())
	if err != nil {
		panic(err)
	}
	client, err := ce.NewClient(sender)
	if err != nil {
		panic(err)
	}
	return &CEMessageBroker{
		cfg:    config,
		log:    log,
		sender: sender,
		client: client,
	}
}

func (broker *CEMessageBroker) Send(ctx context.Context, message ce.Event) error {
	err := broker.client.Send(ctx, message)
	if err != nil {
		if ce.IsUndelivered(err) {
			return fmt.Errorf("message undelivered: %w", err)
		}
	}
	fmt.Println("success")
	return nil
}

func (broker *CEMessageBroker) Close(ctx context.Context) error {
	err := broker.sender.Close(ctx)
	if err != nil {
		return err
	}
	return nil
}
