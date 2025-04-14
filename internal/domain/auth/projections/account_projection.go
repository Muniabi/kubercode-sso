package projections

import (
	"context"
	"encoding/json"
	"log/slog"
	"kubercode-sso/internal/domain/auth/dto"
	"kubercode-sso/internal/domain/auth/events"
	"kubercode-sso/internal/domain/auth/repository"
	"kubercode-sso/internal/domain/auth/values"
	"kubercode-sso/internal/infrastructure/es"

	"github.com/google/uuid"
)

type ProjectionProcessor interface {
	ProcessEvent(ctx context.Context, event es.Event) error
}

type projectionProcessor struct {
	accountRepo repository.AccountRepository
	log         *slog.Logger
}

// NewProjectionProcessor создает новый экземпляр projectionProcessor.
func NewProjectionProcessor(accountRepo repository.AccountRepository, log *slog.Logger) *projectionProcessor {
	return &projectionProcessor{
		accountRepo: accountRepo,
		log:         log,
	}
}

// ProcessEvent обрабатывает события и обновляет проекции.
func (p *projectionProcessor) ProcessEvent(ctx context.Context, event es.Event) error {
	p.log.Info("projectionProcessor")
	p.log.Info("Event type:", slog.String("type", string(event.EventType)))
	switch event.EventType {
	case events.RegisterAccount:
		err := p.handleRegisterAccount(ctx, event)
		if err != nil {
			return err
		}
	case events.ChangePassword:
		err := p.handleChangePassword(ctx, event)
		if err != nil {
			return err
		}
	case events.ChangeEmail:
		err := p.handleChangeEmail(ctx, event)
		if err != nil {
			return err
		}
	case events.RestorePassword:
		err := p.handleRestorePassword(ctx, event)
		if err != nil {
			return err
		}
	default:
		p.log.Warn("Unhandled event type", slog.String("eventType", string(event.EventType)))
	}
	return nil
}

// Обработка события RegisterAccount.
func (p *projectionProcessor) handleRegisterAccount(ctx context.Context, event es.Event) error {
	p.log.Info("projectionProcessor")
	p.log.Info("event data in projectionProcessor", slog.Any("event", event))
	var registerEvent events.RegisterAccountEvent
	err := json.Unmarshal(event.Data, &registerEvent)
	if err != nil {
		return err
	}

	account := dto.UserDTO{
		ID:        registerEvent.Id,
		Email:     registerEvent.Email,
		Password:  registerEvent.Password,
		IsMentor:  registerEvent.IsMentor,
	}
	err = p.accountRepo.Save(ctx, account)
	if err != nil {
		return err
	}
	p.log.Info("Account registered", slog.String("email", registerEvent.Email.ToString()))
	return nil
}

// Обработка события ChangePassword.
func (p *projectionProcessor) handleChangePassword(ctx context.Context, event es.Event) error {
	var changePasswordEvent events.ChangePasswordEvent
	err := json.Unmarshal(event.Data, &changePasswordEvent)
	if err != nil {
		return err
	}

	user, err := p.accountRepo.GetById(ctx, event.AggregateID)
	if err != nil {
		return err
	}
	user.Password = changePasswordEvent.Password
	err = p.accountRepo.Update(ctx, user, user.Email)
	if err != nil {
		return err
	}
	p.log.Info("Password changed", slog.String("email", user.Email.ToString()))
	return nil
}

// Обработка события ChangeEmail.
func (p *projectionProcessor) handleChangeEmail(ctx context.Context, event es.Event) error {
	var changeEmailEvent events.ChangeEmailEvent
	err := json.Unmarshal(event.Data, &changeEmailEvent)
	if err != nil {
		return err
	}

	user, err := p.accountRepo.GetById(ctx, event.AggregateID)
	if err != nil {
		return err
	}
	oldEmail := user.Email
	user.Email = changeEmailEvent.Email
	err = p.accountRepo.Update(ctx, user, oldEmail)
	if err != nil {
		return err
	}
	p.log.Info("Email changed", slog.String("oldEmail", oldEmail.ToString()),
		slog.String("newEmail", changeEmailEvent.Email.ToString()))
	return nil
}

func (p *projectionProcessor) handleRestorePassword(ctx context.Context, event es.Event) error {
	var restorePasswordEvent events.RestorePasswordEvent
	err := json.Unmarshal(event.Data, &restorePasswordEvent)
	if err != nil {
		return err
	}
	user, err := p.accountRepo.GetById(ctx, event.AggregateID)
	if err != nil {
		return err
	}
	user.Password = restorePasswordEvent.Password
	err = p.accountRepo.Update(ctx, user, user.Email)
	if err != nil {
		return err
	}
	p.log.Info("Password changed")
	return nil
}

type AccountProjection struct {
	Id          uuid.UUID
	Email       values.Email
	Password    values.Password
	IsMentor    values.IsMentor
	DeviceToken string
}

func (p *AccountProjection) Handle(event events.Event) error {
	switch e := event.(type) {
	case *events.RegisterAccount:
		p.Id = e.Id
		p.Email = e.Email
		p.Password = e.Password
		p.IsMentor = e.IsMentor
		p.DeviceToken = e.DeviceToken
	}
	return nil
}
