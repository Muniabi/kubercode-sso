package queries

import (
	"golang.org/x/net/context"
	"log/slog"
	"kubercode-sso/config"
	"kubercode-sso/internal/domain/auth/repository"
	"kubercode-sso/internal/infrastructure/es"
)

type GetAllDeviceIdByIdQueryHandler struct {
	es.QueryHandler[GetAllDeviceIdByIdQuery]
	log           *slog.Logger
	cfg           *config.Config
	jwtRepository repository.AbstractTokenRepository
}

func NewGetAllDeviceIdByIdQueryHandler(log *slog.Logger,
	cfg *config.Config, jwtRepository repository.AbstractTokenRepository) *GetAllDeviceIdByIdQueryHandler {
	return &GetAllDeviceIdByIdQueryHandler{
		log:           log,
		cfg:           cfg,
		jwtRepository: jwtRepository,
	}
}

func (handler *GetAllDeviceIdByIdQueryHandler) Handle(ctx context.Context, query es.Query) (interface{}, error) {
	handler.log.Info("GetAllDeviceIdByIdQueryHandler", "query", query)
	getAllDevicesQuery, ok := query.(*GetAllDeviceIdByIdQuery)
	if !ok {
		handler.log.Warn("GetAllDeviceIdByIdQueryHandler", "query", "*getAllDevicesQuery")
	}
	data, err := handler.jwtRepository.GetAllDeviceIdFromEmail(ctx, getAllDevicesQuery.Email)
	if err != nil {
		return nil, err
	}
	return data, err
}
