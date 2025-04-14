package aggregate

import (
	"fmt"
	"kubercode-sso/internal/domain/auth/events"
	"kubercode-sso/internal/domain/auth/values"
	"kubercode-sso/internal/infrastructure/es"

	"github.com/google/uuid"
)

type Account struct {
	es.AggregateRoot
	Email     values.Email
	Password  values.Password
	IsMentor  values.IsMentor
	DeviceToken string
}

func NewAccountWithOnlyId(id uuid.UUID) *Account {
	return &Account{
		AggregateRoot: es.NewAggregateRootWithId(id),
	}
}

func NewAccountWithId(id uuid.UUID, email values.Email, password values.Password, isMentor values.IsMentor) *Account {
	return &Account{
		AggregateRoot: es.NewAggregateRootWithId(id),
		Email:         email,
		Password:      password,
		IsMentor:      isMentor,
	}
}

func (a *Account) When(event es.Event) {
	switch event.EventType {
	case events.RegisterAccount:
		err := a.onCreateAccount(event)
		if err != nil {
			panic(err)
		}
		a.Version++
		fmt.Println(a.Version)
	case events.ChangeEmail:
		err := a.onChangeEmail(event)
		if err != nil {
			panic(err)
		}
		a.Version++
		fmt.Println(a.Version)
	case events.ChangePassword:
		err := a.onChangePassword(event)
		if err != nil {
			panic(err)
		}
		a.Version++
		fmt.Println(a.Version)
	}
}

func (a *Account) onCreateAccount(event es.Event) error {
	var createAccountEvent events.RegisterAccountEvent
	err := event.GetJsonData(&createAccountEvent)
	if err != nil {
		return err
	}
	a.Id = createAccountEvent.Id
	a.Email = createAccountEvent.Email
	a.Password = createAccountEvent.Password
	a.IsMentor = createAccountEvent.IsMentor
	a.DeviceToken = createAccountEvent.DeviceToken
	return nil
}

func (a *Account) onChangePassword(event es.Event) error {
	var changePasswordEvent events.ChangePasswordEvent
	err := event.GetJsonData(&changePasswordEvent)
	if err != nil {
		return err
	}
	a.Id = changePasswordEvent.Id
	a.Password = changePasswordEvent.Password
	return nil
}
func (a *Account) onChangeEmail(event es.Event) error {
	var changeEmailEvent events.ChangeEmailEvent
	err := event.GetJsonData(&changeEmailEvent)
	if err != nil {
		return err
	}
	a.Id = changeEmailEvent.Id
	a.Email = changeEmailEvent.Email
	return nil
}
