package email

type EmailSender interface {
	SendEmail(to string, subject string, body string) error
	SendTestEmail(to string, subject string, body string) error
}
