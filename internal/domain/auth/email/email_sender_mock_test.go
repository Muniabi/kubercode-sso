package email

import (
	"bytes"
	"fmt"
	"net/smtp"
)

// MailHogEmailSender реализует интерфейс EmailSender с использованием MailHog
type MailHogEmailSender struct {
	SMTPServer string
	SMTPPort   string
}

// NewMailHogEmailSender создает новый MailHogEmailSender
func NewMailHogEmailSender() *MailHogEmailSender {
	return &MailHogEmailSender{
		SMTPServer: "localhost",
		SMTPPort:   "1025", // Порт по умолчанию для MailHog
	}
}

// SendEmail отправляет электронное письмо
func (m *MailHogEmailSender) SendEmail(to string, subject string, body string) error {
	return m.send(to, subject, body)
}

// SendTestEmail отправляет тестовое электронное письмо
func (m *MailHogEmailSender) SendTestEmail(to string, subject string, body string) error {
	return m.send(to, subject, body)
}

// send отправляет письмо через SMTP
func (m *MailHogEmailSender) send(to string, subject string, body string) error {
	from := "no-reply@example.com"
	msg := bytes.Buffer{}
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("\r\n")
	msg.WriteString(body)

	err := smtp.SendMail(fmt.Sprintf("%s:%s", m.SMTPServer, m.SMTPPort), nil, from, []string{to}, msg.Bytes())
	if err != nil {
		return err
	}
	return nil
}
