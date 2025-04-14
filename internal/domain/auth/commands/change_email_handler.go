package commands

import (
	"context"
	"fmt"
	"log/slog"
	"kubercode-sso/config"
	aggregate2 "kubercode-sso/internal/domain/auth/aggregate"
	"kubercode-sso/internal/domain/auth/events"
	"kubercode-sso/internal/infrastructure/es"
	"kubercode-sso/internal/infrastructure/es/store"
)

type ChangeEmailHandler struct {
	es.CommandHandler[ChangeEmailCommand]
	log            *slog.Logger
	cfg            *config.Config
	eventStore     store.EventStore
	aggregateStore store.AggregateStore
}

func NewChangeEmailHandler(log *slog.Logger, cfg *config.Config, eventStore store.EventStore,
	aggregateStore store.AggregateStore) *ChangeEmailHandler {
	return &ChangeEmailHandler{log: log, cfg: cfg, eventStore: eventStore, aggregateStore: aggregateStore}
}

func (c *ChangeEmailHandler) Handle(ctx context.Context, command es.Command) (es.Event, error) {
	c.log.Info("ChangePasswordHandler", "handle")
	if command == nil {
		return es.Event{}, fmt.Errorf("received nil command")
	}

	changeEmailCommand, ok := command.(*ChangeEmailCommand)
	if !ok {
		return es.Event{}, fmt.Errorf("invalid command type: expected *CreateAccountCommand, got %T", command)
	}
	aggregate := aggregate2.NewAccountWithOnlyId(changeEmailCommand.AggregateID)
	_ = c.aggregateStore.SaveEventsToAggregate(ctx, aggregate)
	_ = c.aggregateStore.LoadAndApplyEvents(ctx, aggregate)
	c.log.Info("accountId", aggregate.Email)
	event, err := events.NewChangeEmailEvent(aggregate.Id, changeEmailCommand.NewEmail, aggregate)
	if err != nil {
		return es.Event{}, err
	}
	err = c.eventStore.SaveEvents(ctx, aggregate.Id, event)
	if err != nil {
		return es.Event{}, err
	}
	return event, nil
}
