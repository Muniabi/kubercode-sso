package es

import (
	"encoding/json"
	"fmt"
	"github.com/EventStore/EventStore-Client-Go/v4/esdb"
	"github.com/google/uuid"
	"time"
)

type IEventable interface {
	GetEventID() uuid.UUID
	GetTimeStamp() time.Time
	GetData() []byte
	SetData([]byte) *Event
	GetJsonData(data interface{}) error
	GetEventType() EventType
	GetAggregateType() AggregateType
	SetAggregateType(aggregateType AggregateType)
}

// EventType тип любого event, используется как уникальный id.
type EventType string

type Event struct {
	/*
		Event - это внутреннее представление события, возвращаемое, когда Aggregate использует
		NewEvent для создания нового события. События, загружаемые из базы данных,
		представлены каждым внутренним типом события DBs, реализующим Event.
	*/
	IEventable
	EventID       uuid.UUID
	EventType     EventType
	Data          []byte // json.Marshall
	Timestamp     time.Time
	AggregateType AggregateType
	AggregateID   uuid.UUID
	Version       int
	Metadata      []byte
}

func NewBaseEvent(aggregate IAggregateRoot, eventType EventType) Event {
	/* Конструктор Event */
	return Event{
		EventID:       uuid.New(),
		AggregateType: aggregate.GetType(),
		AggregateID:   aggregate.GetId(),
		Version:       aggregate.GetVersion(),
		EventType:     eventType,
		Timestamp:     time.Now().UTC(),
	}
}

func NewEventFromRecorded(event *esdb.RecordedEvent) (Event, error) {
	aggregateId, err := uuid.Parse(event.StreamID[5:])
	if err != nil {
		return Event{}, err
	}
	return Event{
		EventID:     event.EventID,
		EventType:   EventType(event.EventType),
		Data:        event.Data,
		Timestamp:   event.CreatedDate,
		AggregateID: aggregateId,
		Version:     int(event.EventNumber),
		Metadata:    event.UserMetadata,
	}, nil
}

func (e *Event) GetEventID() uuid.UUID {
	return e.EventID
}

func (e *Event) GetTimeStamp() time.Time {
	return e.Timestamp
}

func (e *Event) GetData() []byte {
	return e.Data
}

func (e *Event) SetData(data []byte) *Event {
	e.Data = data
	return e
}

func (e *Event) GetJsonData(data interface{}) error {
	return json.Unmarshal(e.GetData(), data)
}

func (e *Event) SetJsonData(data interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	e.Data = dataBytes
	return nil
}

func (e *Event) GetEventType() EventType {
	return e.EventType
}

func (e *Event) GetAggregateType() AggregateType {
	return e.AggregateType
}

func (e *Event) SetAggregateType(aggregateType AggregateType) {
	e.AggregateType = aggregateType
}

func (e *Event) GetAggregateID() uuid.UUID {
	return e.AggregateID
}

func (e *Event) GetVersion() int {
	return e.Version
}

func (e *Event) SetVersion(aggregateVersion int) {
	e.Version = aggregateVersion
}

func (e *Event) GetMetadata() []byte {
	return e.Metadata
}

func (e *Event) SetMetadata(metaData interface{}) error {

	metaDataBytes, err := json.Marshal(metaData)
	if err != nil {
		return err
	}

	e.Metadata = metaDataBytes
	return nil
}

func (e *Event) GetJsonMetadata(metaData interface{}) error {
	return json.Unmarshal(e.GetMetadata(), metaData)
}

func (e *Event) GetString() string {
	return fmt.Sprintf("event: %+v", e)
}

func (e *Event) String() string {
	return fmt.Sprintf("(Event): AggregateID: {%s}, Version: {%d}, EventType: {%s}, AggregateType: {%s}, Metadata: {%s}, TimeStamp: {%s}",
		e.AggregateID,
		e.Version,
		e.EventType,
		e.AggregateType,
		string(e.Metadata),
		e.Timestamp.UTC().String(),
	)
}
