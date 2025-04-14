package events

import (
	"kubercode-sso/internal/domain/auth/values"
	"kubercode-sso/internal/infrastructure/es"
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
)

const (
	RegisterAccount = "RegisterAccount"
	ChangeEmail     = "ChangeEmail"
	ChangePassword  = "ChangePassword"
	RestorePassword = "RestorePassword"
	SendEmail       = "SendEmail"
)

type RegisterAccountEvent struct {
	Id        uuid.UUID        `json:"id"`
	Email     values.Email     `json:"email"`
	Password  values.Password  `json:"password"`
	IsMentor  values.IsMentor   `json:"is_mentor"`
	DeviceToken string         `json:"device_token"`
}

func NewRegisterAccount(id uuid.UUID, email values.Email, password values.Password, isMentor values.IsMentor,
	deviceToken string, aggregate es.IAggregateRoot) (es.Event, error) {
	registerAccount := &RegisterAccountEvent{
		Id:        id,
		Email:     email,
		Password:  password,
		IsMentor:  isMentor,
		DeviceToken: deviceToken,
	}
	event := es.NewBaseEvent(aggregate, RegisterAccount)
	err := event.SetJsonData(&registerAccount)
	if err != nil {
		return es.Event{}, err
	}
	return event, nil
}

type ChangeEmailEvent struct {
	Id    uuid.UUID    `json:"id"`
	Email values.Email `json:"email"`
}

func NewChangeEmailEvent(id uuid.UUID, email values.Email, aggregate es.IAggregateRoot) (es.Event, error) {
	changeEmail := &ChangeEmailEvent{
		Id:    id,
		Email: email,
	}
	event := es.NewBaseEvent(aggregate, ChangeEmail)
	err := event.SetJsonData(&changeEmail)
	if err != nil {
		return es.Event{}, err
	}
	return event, nil
}

type ChangePasswordEvent struct {
	Id       uuid.UUID       `json:"id"`
	Password values.Password `json:"password"`
}

func NewChangePasswordEvent(id uuid.UUID, password values.Password, aggregate es.IAggregateRoot) (es.Event, error) {
	changePassword := &ChangePasswordEvent{
		Id:       id,
		Password: password,
	}
	event := es.NewBaseEvent(aggregate, ChangePassword)
	err := event.SetJsonData(&changePassword)
	if err != nil {
		return es.Event{}, err
	}
	return event, nil
}

type RestorePasswordEvent struct {
	Id       uuid.UUID       `json:"id"`
	Password values.Password `json:"password"`
}

func NewRestorePasswordEvent(id uuid.UUID, password values.Password, aggregate es.IAggregateRoot) (es.Event, error) {
	restorePassword := &RestorePasswordEvent{
		Id:       id,
		Password: password,
	}
	event := es.NewBaseEvent(aggregate, RestorePassword)
	err := event.SetJsonData(&restorePassword)
	if err != nil {
		return es.Event{}, err
	}
	return event, nil
}

type SendEmailEvent struct {
	Email   values.Email `json:"email"`
	Subject string       `json:"subject"`
	Body    string       `json:"body"`
}

func NewSendEmailEvent(email values.Email, subject string, body string, aggregate es.IAggregateRoot) (es.Event, error) {
	sendEmailEvent := &SendEmailEvent{
		Email:   email,
		Subject: subject,
		Body:    body,
	}
	event := es.NewBaseEvent(aggregate, SendEmail)
	err := event.SetJsonData(sendEmailEvent)
	if err != nil {
		return es.Event{}, err
	}
	return event, nil
}

func NewSendEmailEventToJson(email values.Email, subject string, body string) *SendEmailEvent {
	return &SendEmailEvent{
		Email:   email,
		Subject: subject,
		Body:    body,
	}
}

func (event *SendEmailEvent) ToCloudEvent(id, eventType string, data *SendEmailEvent) (ce.Event, error) {
	cloudevent := ce.NewEvent("1.0")
	cloudevent.SetID(id)
	cloudevent.SetType(eventType)
	cloudevent.SetTime(time.Now())
	cloudevent.SetSource("kubercode-sso")
	cloudevent.SetSpecVersion("1.0")
	_ = cloudevent.SetData(ce.ApplicationJSON, data)
	return cloudevent, nil

}
