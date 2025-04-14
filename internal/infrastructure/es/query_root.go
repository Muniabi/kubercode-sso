package es

import (
	"context"
	"github.com/google/uuid"
)

type Query interface {
	GetAggregateID() uuid.UUID
}

type QueryResponse string // временный тип, который в себе содержит ответ из read-модели

type QueryHandler[AbstractQuery Query] interface {
	Handle(ctx context.Context, query AbstractQuery) (interface{}, error)
}

type BaseQuery struct {
	AggregateID uuid.UUID
}

func NewBaseQuery(aggregateID uuid.UUID) BaseQuery {
	return BaseQuery{
		AggregateID: aggregateID,
	}
}

func (q BaseQuery) GetAggregateID() uuid.UUID {
	return q.AggregateID
}
