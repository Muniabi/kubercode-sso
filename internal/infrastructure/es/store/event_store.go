package store

import (
	"context"
	"errors"
	"fmt"
	"io"
	"kubercode-sso/config"
	"kubercode-sso/internal/domain/auth/projections"
	"kubercode-sso/internal/infrastructure/es"
	"log/slog"

	"github.com/EventStore/EventStore-Client-Go/v4/esdb"
	"github.com/google/uuid"
)

type EventStore interface {
	SaveEvents(ctx context.Context, aggregateId uuid.UUID, events ...es.Event) error
	LoadEvents(ctx context.Context, aggregateId uuid.UUID) ([]es.Event, error)
	CreatePersistentSub()
	SubscribeToStream(ctx context.Context) error
	EventDataFromEvent(event es.Event) esdb.EventData
	EventFromData(recordedEvent *esdb.RecordedEvent) (es.Event, error)
}

// eventStoreDB представляет собой реализацию EventStore для eventStoreDB.
type eventStoreDB struct {
	cfg       *config.Config
	log       *slog.Logger
	Client    *esdb.Client
	processor projections.ProjectionProcessor
}

// NewEventStore создает новый экземпляр eventStoreDB.
func NewEventStore(cfg *config.Config, log *slog.Logger, processor projections.ProjectionProcessor) (*eventStoreDB, error) {
	settings, err := esdb.ParseConnectionString(cfg.EventStoreConnectionString)
	if err != nil {
		panic(err)
	}
	db, err := esdb.NewClient(settings)
	if err != nil {
		panic(err)
	}
	return &eventStoreDB{
		cfg:       cfg,
		log:       log,
		Client:    db,
		processor: processor,
	}, nil
}

// SaveEvents сохраняет события для указанного aggregateId.
func (esd *eventStoreDB) SaveEvents(ctx context.Context, aggregateId uuid.UUID, events ...es.Event) error {
	streamName := fmt.Sprintf("user-%s", aggregateId.String())
	esd.log.Info("saving events", slog.String("streamName", streamName))
	var eventData []esdb.EventData
	for _, event := range events {
		eventData = append(eventData, esd.EventDataFromEvent(event))
	}
	esd.log.Info("eventData", eventData)
	_, err := esd.Client.AppendToStream(ctx, streamName, esdb.AppendToStreamOptions{}, eventData...)
	return err
}

// LoadEvents загружает все события для указанного aggregateId.
func (esd *eventStoreDB) LoadEvents(ctx context.Context, aggregateId uuid.UUID) ([]es.Event, error) {
	streamName := fmt.Sprintf("user-%s", aggregateId.String())
	esd.log.Info("loading events", slog.String("streamName", streamName))
	var events []es.Event
	options := esdb.ReadStreamOptions{
		From:      esdb.Start{},
		Direction: esdb.Forwards,
	}
	res, err := esd.Client.ReadStream(ctx, streamName, options, 100)
	defer res.Close()
	if err != nil {
		return nil, err
	}
	for {
		event, err := res.Recv()

		if errors.Is(err, io.EOF) {
			break
		}
		myEvent, err := esd.EventFromData(event.OriginalEvent())
		if err != nil {
			return nil, err
		}
		events = append(events, myEvent)
	}
	return events, nil
}

func (esd *eventStoreDB) CreatePersistentSub() {
	options := esdb.PersistentAllSubscriptionOptions{
		Filter: &esdb.SubscriptionFilter{
			Type:     esdb.StreamFilterType,
			Prefixes: []string{esd.cfg.AccountPrefix},
		},
	}

	err := esd.Client.CreatePersistentSubscriptionToAll(context.Background(), esd.cfg.ProjectionsGroupName, options)

	if err != nil {
		esd.log.Warn("Failed to create persistent subscription", "err", err)
		return
	}
}

// SubscribeToStream подписывается на поток событий.
func (esd *eventStoreDB) SubscribeToStream(ctx context.Context) error {
	esd.log.Info("SubscribeToStream")
	sub, err := esd.Client.SubscribeToPersistentSubscriptionToAll(context.Background(),
		esd.cfg.ProjectionsGroupName, esdb.SubscribeToPersistentSubscriptionOptions{})

	if err != nil {
		panic(err)
	}
	eventChannel := make(chan es.Event, 10)
	defer close(eventChannel)
	go func() {
		for {
			event := sub.Recv()

			if event.EventAppeared != nil {
				sub.Ack(event.EventAppeared.Event)
				esd.log.Info("MY EVENTS LMAO QEQOQEQ", slog.Any("aaa", event.EventAppeared.Event.Event))
				myEvent, err := esd.EventFromData(event.EventAppeared.Event.Event)
				if err != nil {
					esd.log.Error("Failed to parse event", "err", err)
					close(eventChannel)
				}
				err = esd.processor.ProcessEvent(ctx, myEvent)
				if err != nil {
					esd.log.Error("Error while handle event", slog.Any("error", err))
					close(eventChannel)
				}
			}

			if event.SubscriptionDropped != nil {
				break
			}
		}
	}()
	return nil
}

// EventDataFromEvent преобразует пользовательское событие в eventStoreDB EventData.
func (esd *eventStoreDB) EventDataFromEvent(event es.Event) esdb.EventData {
	return esdb.EventData{
		EventType: string(event.EventType),
		Data:      event.Data,
		EventID:   event.EventID,
	}
}

// EventFromData преобразует EventData из eventStoreDB в пользовательское событие.
func (esd *eventStoreDB) EventFromData(recordedEvent *esdb.RecordedEvent) (es.Event, error) {
	var event es.Event
	event, err := es.NewEventFromRecorded(recordedEvent)
	if err != nil {
		return event, err
	}
	return event, nil
}
