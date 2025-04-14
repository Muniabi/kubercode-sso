package fabric

import (
	"fmt"
	"reflect"
	"kubercode-sso/internal/infrastructure/es"
)

type HandlerFabric struct {
	commandHandlers map[reflect.Type]es.CommandHandler[es.Command]
	queryHandlers   map[reflect.Type]es.QueryHandler[es.Query]
}

func NewHandlerFabric() *HandlerFabric {
	return &HandlerFabric{
		commandHandlers: make(map[reflect.Type]es.CommandHandler[es.Command]),
		queryHandlers:   make(map[reflect.Type]es.QueryHandler[es.Query]),
	}
}

func (handlerFabric *HandlerFabric) RegisterCommandHandler(command es.Command, handler es.CommandHandler[es.Command]) {
	handlerFabric.commandHandlers[reflect.TypeOf(command)] = handler
}
func (handlerFabric *HandlerFabric) RegisterQueryHandler(query es.Query, handler es.QueryHandler[es.Query]) {
	handlerFabric.queryHandlers[reflect.TypeOf(query)] = handler
}
func (handlerFabric *HandlerFabric) GetCommandHandler(command es.Command) (es.CommandHandler[es.Command], error) {
	handler, ok := handlerFabric.commandHandlers[reflect.TypeOf(command)]
	if !ok {
		return nil, fmt.Errorf("command handler not registered")
	}
	return handler, nil
}
func (handlerFabric *HandlerFabric) GetQueryHandler(query es.Query) (es.QueryHandler[es.Query], error) {
	handler, ok := handlerFabric.queryHandlers[reflect.TypeOf(query)]
	if !ok {
		return nil, fmt.Errorf("query handler not registered")
	}
	return handler, nil
}
