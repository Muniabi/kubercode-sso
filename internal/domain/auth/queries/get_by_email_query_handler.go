package queries

import (
	"context"
	"log/slog"
	"kubercode-sso/config"
	"kubercode-sso/internal/domain/auth/repository"
	"kubercode-sso/internal/infrastructure/es"
)

type GetByEmailQueryHandler struct {
	es.QueryHandler[GetByEmailQuery]
	log         *slog.Logger
	cfg         *config.Config
	accountRepo repository.AccountRepository
}

func NewGetByEmailQueryHandler(accountRepo repository.AccountRepository, log *slog.Logger, cfg *config.Config) *GetByEmailQueryHandler {
	return &GetByEmailQueryHandler{
		log:         log,
		accountRepo: accountRepo,
		cfg:         cfg,
	}
}

func (handler *GetByEmailQueryHandler) Handle(ctx context.Context, query es.Query) (interface{}, error) {
	handler.log.Info("GetByEmailQueryHandler", "query", query)
	verifyDataQuery, ok := query.(*GetByEmailQuery)
	if !ok {
		handler.log.Warn("GetByEmailQueryHandler", "query", "*GetByEmailQuery")
	}
	data, err := handler.accountRepo.GetByEmail(ctx, verifyDataQuery.Email)
	if err != nil {
		handler.log.Warn("GetByEmailQueryHandler", "query", "GetByEmail", "error", err)
		return nil, err
	}
	return data, nil
}
