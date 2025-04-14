package store

import (
	"context"
	"log/slog"
	"kubercode-sso/internal/infrastructure/es"
)

type AggregateStore interface {
	SaveEventsToAggregate(ctx context.Context, aggregate es.IAggregateRoot) error
	LoadAndApplyEvents(ctx context.Context, aggregate es.IAggregateRoot) error
}

type EsAggregateStore struct {
	es  EventStore
	log *slog.Logger
}

func NewEsAggregateStore(eventStore EventStore, logger *slog.Logger) *EsAggregateStore {
	return &EsAggregateStore{
		es:  eventStore,
		log: logger,
	}
}

func (esa *EsAggregateStore) SaveEventsToAggregate(ctx context.Context, aggregate es.IAggregateRoot) error {
	// сохраняет все события(аппендит к списочку внутри агрегата)
	esa.log.Info("Saving events to aggregate", aggregate)
	events, err := esa.es.LoadEvents(ctx, aggregate.GetId())
	if err != nil {
		return err
	}
	aggregate.AddEvents(events...)
	return nil
}

func (esa *EsAggregateStore) LoadAndApplyEvents(ctx context.Context, aggregate es.IAggregateRoot) error {
	// тут применяются все события и возвращается агрегат

	allEventsFromAggregate := aggregate.GetEvents()
	for _, event := range allEventsFromAggregate {
		aggregate.When(event)
	}
	aggregate.ClearEvents()
	return nil
}
