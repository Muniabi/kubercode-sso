package email

import (
	"testing"
)

func TestMailHogEmailSender(t *testing.T) {
	sender := NewMailHogEmailSender()
	to := "test@example.com"
	subject := "Test Email"
	body := "This is a test email."

	// Отправка обычного email
	err := sender.SendEmail(to, subject, body)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Отправка тестового email
	err = sender.SendTestEmail(to, subject, body)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
