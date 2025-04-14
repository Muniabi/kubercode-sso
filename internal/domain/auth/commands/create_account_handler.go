package commands

import (
	"context"
	"fmt"
	"log/slog"
	"kubercode-sso/config"
	auth "kubercode-sso/internal/domain/auth/aggregate"
	"kubercode-sso/internal/domain/auth/events"
	"kubercode-sso/internal/infrastructure/es"
	"kubercode-sso/internal/infrastructure/es/store"
)

type CreateAccountHandler struct {
	es.CommandHandler[CreateAccountCommand]
	log        *slog.Logger
	cfg        *config.Config
	eventStore store.EventStore
}

func NewCreateAccountHandler(log *slog.Logger, cfg *config.Config, eventStore store.EventStore) *CreateAccountHandler {
	return &CreateAccountHandler{log: log, cfg: cfg, eventStore: eventStore}
}

func (c *CreateAccountHandler) Handle(ctx context.Context, command es.Command) (es.Event, error) {
	c.log.Info("CreateAccountHandler", "handle")
	if command == nil {
		return es.Event{}, fmt.Errorf("received nil command")
	}

	createAccountCommand, ok := command.(*CreateAccountCommand)
	if !ok {
		return es.Event{}, fmt.Errorf("invalid command type: expected *CreateAccountCommand, got %T", command)
	}

	aggregate := auth.NewAccountWithId(createAccountCommand.AggregateID, createAccountCommand.Email,
		createAccountCommand.Password, createAccountCommand.IsCompany)
	event, err := events.NewRegisterAccount(createAccountCommand.AggregateID, createAccountCommand.Email,
		createAccountCommand.Password, createAccountCommand.IsCompany, aggregate)
	if err != nil {
		return es.Event{}, err
	}
	err = c.eventStore.SaveEvents(ctx, aggregate.Id, event)
	if err != nil {
		return es.Event{}, err
	}
	return event, nil
}
