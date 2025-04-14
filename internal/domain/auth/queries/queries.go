package queries

import (
	"github.com/google/uuid"
	"kubercode-sso/internal/domain/auth/values"
	"kubercode-sso/internal/infrastructure/es"
)

type GetByEmailQuery struct {
	es.BaseQuery
	Email values.Email
}

func NewEmptyGetByEmailQuery() *GetByEmailQuery {
	return &GetByEmailQuery{}
}

func NewGetByEmailQuery(aggregateId uuid.UUID, email values.Email) *GetByEmailQuery {
	return &GetByEmailQuery{
		BaseQuery: es.NewBaseQuery(aggregateId),
		Email:     email,
	}
}

type GetAllDeviceIdByIdQuery struct {
	es.BaseQuery
	Email string `bson:"user_email"`
}

func NewEmptyGetAllDeviceIdByIdQuery() *GetAllDeviceIdByIdQuery {
	return &GetAllDeviceIdByIdQuery{}
}

func NewGetAllDeviceIdByIdQuery(aggregateId uuid.UUID, email string) *GetAllDeviceIdByIdQuery {
	return &GetAllDeviceIdByIdQuery{
		BaseQuery: es.NewBaseQuery(aggregateId),
		Email:     email,
	}
}

type GetUserByIdQuery struct {
	es.BaseQuery
	Id uuid.UUID `bson:"_id"`
}

func NewEmptyGetUserByIdQuery() *GetUserByIdQuery {
	return &GetUserByIdQuery{}
}

func NewGetUserByIdQuery(aggregateId uuid.UUID, id uuid.UUID) *GetUserByIdQuery {
	return &GetUserByIdQuery{
		BaseQuery: es.NewBaseQuery(aggregateId),
		Id:        id,
	}
}
