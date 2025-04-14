package commands

import (
	"github.com/google/uuid"
	"kubercode-sso/internal/domain/auth/values"
	"kubercode-sso/internal/infrastructure/es"
)

type CreateAccountCommand struct {
	es.BaseCommand
	Email     values.Email
	Password  values.Password
	IsCompany values.IsCompany
}

func NewCreateEmptyAccountCommand() *CreateAccountCommand {
	return &CreateAccountCommand{}
}

func NewCreateAccountCommand(aggregateID uuid.UUID, email values.Email, password values.Password,
	isCompany values.IsCompany) *CreateAccountCommand {
	return &CreateAccountCommand{BaseCommand: es.NewBaseCommand(aggregateID), Email: email, Password: password,
		IsCompany: isCompany}
}

type ChangePasswordCommand struct {
	es.BaseCommand
	OldPassword values.Password
	NewPassword values.Password
}

func NewEmptyChangePasswordCommand() *ChangePasswordCommand {
	return &ChangePasswordCommand{}
}

func NewChangePasswordCommand(aggregateID uuid.UUID, oldPassword values.Password,
	newPassword values.Password) *ChangePasswordCommand {
	return &ChangePasswordCommand{
		BaseCommand: es.NewBaseCommand(aggregateID),
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}
}

type ChangeEmailCommand struct {
	es.BaseCommand
	NewEmail values.Email
}

func NewEmptyChangeEmailCommand() *ChangeEmailCommand {
	return &ChangeEmailCommand{}
}

func NewChangeEmailCommand(aggregateID uuid.UUID,
	newEmail values.Email) *ChangeEmailCommand {
	return &ChangeEmailCommand{
		BaseCommand: es.NewBaseCommand(aggregateID),
		NewEmail:    newEmail,
	}
}

type RestorePasswordCommand struct {
	es.BaseCommand
	Password values.Password
}

func NewEmptyRestorePasswordCommand() *RestorePasswordCommand {
	return &RestorePasswordCommand{}
}

func NewRestorePasswordCommand(aggregateID uuid.UUID, password values.Password) *RestorePasswordCommand {
	return &RestorePasswordCommand{
		BaseCommand: es.NewBaseCommand(aggregateID),
		Password:    password,
	}
}
