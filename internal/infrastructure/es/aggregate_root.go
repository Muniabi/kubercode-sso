package es

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

const (
	aggregateStartVersion = -1
)

// тип агрегата

type AggregateType string

// базовый интерфейс с описанием всех методов AggregateRoot

type IAggregateRoot interface {
	GetId() uuid.UUID
	GetVersion() int
	GetType() AggregateType
	SetType(aggregateType AggregateType)
	GetCreationTime() time.Time
	GetUpdateTime() time.Time
	GetEvents() []Event
	When(Event)
	AddEvent(event Event)
	AddEvents(events ...Event)
	ClearEvents()
}

// структура, которая представляет собой базовый агрегат с всеми нужными полями

type AggregateRoot struct {
	IAggregateRoot
	Id           uuid.UUID
	Version      int
	Type         AggregateType
	CreationTime time.Time
	UpdateTime   time.Time
	Events       []Event
}

func NewAggregateRoot() AggregateRoot {
	return AggregateRoot{
		Id:           uuid.New(),
		Version:      aggregateStartVersion,
		CreationTime: time.Now(),
		UpdateTime:   time.Now(),
		Events:       make([]Event, 0),
	}
}

func NewAggregateRootWithId(id uuid.UUID) AggregateRoot {
	if id == uuid.Nil {
		return AggregateRoot{}
	}
	return AggregateRoot{
		Id:           id,
		Version:      aggregateStartVersion,
		CreationTime: time.Now(),
		UpdateTime:   time.Now(),
		Events:       make([]Event, 0),
	}
}

func (a *AggregateRoot) When(event Event) {}

func (a *AggregateRoot) GetId() uuid.UUID {
	return a.Id
}

func (a *AggregateRoot) GetVersion() int {
	return a.Version
}

func (a *AggregateRoot) GetType() AggregateType {
	return a.Type
}

func (a *AggregateRoot) SetType(aggregateType AggregateType) {
	a.Type = aggregateType
}

func (a *AggregateRoot) GetCreationTime() time.Time {
	return a.CreationTime
}

func (a *AggregateRoot) GetUpdateTime() time.Time {
	return a.UpdateTime
}
func (a *AggregateRoot) GetEvents() []Event {
	return a.Events
}
func (a *AggregateRoot) AddEvent(event Event) {
	a.Events = append(a.Events, event)
}
func (a *AggregateRoot) AddEvents(events ...Event) {
	for _, event := range events {
		a.AddEvent(event)
	}
}
func (a *AggregateRoot) ClearEvents() {
	a.Events = make([]Event, 0)
}

func (a *AggregateRoot) String() string {
	return fmt.Sprintf("(Aggregate): Id:%s, Version:%s, Type:%s, CreationTime:%s, UpdateTime:%s, Events:%v",
		a.Id, a.Version, a.Type, a.CreationTime, a.UpdateTime, len(a.Events))
}
