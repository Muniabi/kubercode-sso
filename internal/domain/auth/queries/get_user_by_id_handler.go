package queries

import (
	"context"
	"log/slog"
	"kubercode-sso/config"
	"kubercode-sso/internal/domain/auth/repository"
	"kubercode-sso/internal/infrastructure/es"
)

type GetUserByIDQueryHandler struct {
	es.QueryHandler[GetUserByIdQuery]
	log         *slog.Logger
	cfg         *config.Config
	accountRepo repository.AccountRepository
}

func NewGetUserByIDQueryHandler(log *slog.Logger, cfg *config.Config, accountRepo repository.AccountRepository) *GetUserByIDQueryHandler {
	return &GetUserByIDQueryHandler{
		log:         log,
		cfg:         cfg,
		accountRepo: accountRepo,
	}
}

func (handler *GetUserByIDQueryHandler) Handle(ctx context.Context, query es.Query) (interface{}, error) {
	handler.log.Info("GetUserByIDQueryHandler", "query", query)
	getByIdQuery, ok := query.(*GetUserByIdQuery)
	if !ok {
		handler.log.Warn("GetUserByIDQueryHandler", "query", "*GetByEmailQuery")
	}
	data, err := handler.accountRepo.GetById(ctx, getByIdQuery.Id)
	if err != nil {
		handler.log.Warn("GetUserByIDQueryHandler", "query", "GetById", "err", err)
		return nil, err
	}
	return data, nil
}
