package es

import (
	"context"
	"github.com/google/uuid"
)

type Command interface {
	GetAggregateID() uuid.UUID
}

type CommandHandler[AbstractCommand Command] interface {
	Handle(ctx context.Context, command AbstractCommand) (Event, error)
}

type BaseCommand struct {
	AggregateID uuid.UUID
}

func NewBaseCommand(aggregateID uuid.UUID) BaseCommand {
	return BaseCommand{AggregateID: aggregateID}
}

func (c BaseCommand) GetAggregateID() uuid.UUID {
	return c.AggregateID
}
