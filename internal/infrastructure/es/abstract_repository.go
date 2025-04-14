package es

import (
	"context"
	"github.com/google/uuid"
	sso "kubercode-sso/proto/pb/go"
)

type AbstractRepository interface {
	Get(ctx context.Context, accountId uuid.UUID) (string, error)
	Update(ctx context.Context, data sso.Account) error
	Create(ctx context.Context, account sso.Account) error
}
