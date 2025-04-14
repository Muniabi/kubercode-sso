package commands

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"kubercode-sso/config"
	aggregate2 "kubercode-sso/internal/domain/auth/aggregate"
	"kubercode-sso/internal/domain/auth/events"
	"kubercode-sso/internal/infrastructure/es"
	"kubercode-sso/internal/infrastructure/es/store"
)

type ChangePasswordHandler struct {
	es.CommandHandler[ChangePasswordCommand]
	log            *slog.Logger
	cfg            *config.Config
	eventStore     store.EventStore
	aggregateStore store.AggregateStore
}

func NewChangePasswordHandler(log *slog.Logger, cfg *config.Config, eventStore store.EventStore,
	aggregateStore store.AggregateStore) *ChangePasswordHandler {
	return &ChangePasswordHandler{log: log, cfg: cfg, eventStore: eventStore, aggregateStore: aggregateStore}
}

func (c *ChangePasswordHandler) Handle(ctx context.Context, command es.Command) (es.Event, error) {
	c.log.Info("ChangePasswordHandler", "handle")
	if command == nil {
		return es.Event{}, fmt.Errorf("received nil command")
	}

	changePasswordCommand, ok := command.(*ChangePasswordCommand)
	if !ok {
		return es.Event{}, fmt.Errorf("invalid command type: expected *CreateAccountCommand, got %T", command)
	}
	aggregate := aggregate2.NewAccountWithOnlyId(changePasswordCommand.AggregateID)
	_ = c.aggregateStore.SaveEventsToAggregate(ctx, aggregate)
	_ = c.aggregateStore.LoadAndApplyEvents(ctx, aggregate)
	c.log.Info("accountId", aggregate.Email)
	fmt.Println(changePasswordCommand.OldPassword)
	fmt.Println(aggregate.Password)
	if bytes.Equal(changePasswordCommand.OldPassword.GetPassword(), aggregate.Password.GetPassword()) {
		return es.Event{}, fmt.Errorf("paswords dont match")
	}
	event, err := events.NewChangePasswordEvent(aggregate.Id, changePasswordCommand.NewPassword, aggregate)
	if err != nil {
		return es.Event{}, err
	}
	err = c.eventStore.SaveEvents(ctx, aggregate.Id, event)
	if err != nil {
		return es.Event{}, err
	}
	return event, nil
}
