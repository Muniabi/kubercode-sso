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

type RestorePasswordHandler struct {
	es.CommandHandler[RestorePasswordCommand]
	log            *slog.Logger
	cfg            *config.Config
	eventStore     store.EventStore
	aggregateStore store.AggregateStore
}

func NewRestorePasswordHandler(log *slog.Logger, cfg *config.Config, eventStore store.EventStore,
	aggregateStore store.AggregateStore) *RestorePasswordHandler {
	return &RestorePasswordHandler{
		log:            log,
		cfg:            cfg,
		eventStore:     eventStore,
		aggregateStore: aggregateStore,
	}
}

func (handler *RestorePasswordHandler) Handle(ctx context.Context, command es.Command) (es.Event, error) {
	handler.log.Info("RestorePasswordHandler")
	if command == nil {
		return es.Event{}, fmt.Errorf("received nil command")
	}
	restorePasswordCommand, ok := command.(*RestorePasswordCommand)
	if !ok {
		return es.Event{}, fmt.Errorf("received invalid command")
	}
	aggregate := aggregate2.NewAccountWithOnlyId(restorePasswordCommand.AggregateID)
	_ = handler.aggregateStore.SaveEventsToAggregate(ctx, aggregate)
	_ = handler.aggregateStore.LoadAndApplyEvents(ctx, aggregate)
	handler.log.Info("accountId", aggregate.Email)
	event, err := events.NewRestorePasswordEvent(aggregate.Id, restorePasswordCommand.Password, aggregate)
	if err != nil {
		return es.Event{}, err
	}
	err = handler.eventStore.SaveEvents(ctx, aggregate.Id, event)
	if err != nil {
		return es.Event{}, err
	}
	return event, nil

}
